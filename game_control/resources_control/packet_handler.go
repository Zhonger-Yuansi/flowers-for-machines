package resources_control

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/py_rpc"
	"github.com/pterm/pterm"
)

// command request callback
func (r *Resources) handleCommandOutput(p *packet.CommandOutput) {
	callback, ok := r.commandCallback.LoadAndDelete(p.CommandOrigin.UUID)
	if ok {
		callback(p)
	}
}

// heart beat response (netease pyrpc)
func (r *Resources) handlePyRpc(p *packet.PyRpc) {
	// prepare
	if p.Value == nil {
		return
	}
	// unmarshal
	content, err := py_rpc.Unmarshal(p.Value)
	if err != nil {
		pterm.Warning.Sprintf("handlePyRpc: %v", err)
		return
	}
	// unmarshal
	switch c := content.(type) {
	case *py_rpc.HeartBeat:
		// heart beat to test the device is still alive?
		// it seems that we just need to return it back to the server is OK
		c.Type = py_rpc.ClientToServerHeartBeat
		r.client.Conn().WritePacket(&packet.PyRpc{
			Value:         py_rpc.Marshal(c),
			OperationType: packet.PyRpcOperationTypeSend,
		})
	}
}

// inventory contents(basic)
func (r *Resources) handleInventoryContent(p *packet.InventoryContent) {
	windowID := WindowID(p.WindowID)
	for key, value := range p.Content {
		r.inventory.setItemStack(windowID, SlotID(key), &value)
		callbacks, ok := r.inventoryCallback.LoadAndDelete(SlotLocation{WindowID: windowID, SlotID: SlotID(key)})
		if ok {
			callbacks.FinishAll()
		}
	}
}

// inventory contents(for enchant command...)
func (r *Resources) handleInventoryTransaction(p *packet.InventoryTransaction) {
	for _, value := range p.Actions {
		if value.SourceType == protocol.InventoryActionSourceCreative {
			continue
		}

		windowID, slotID := WindowID(value.WindowID), SlotID(value.InventorySlot)
		r.inventory.setItemStack(windowID, slotID, &value.NewItem)

		callbacks, ok := r.inventoryCallback.LoadAndDelete(SlotLocation{WindowID: windowID, SlotID: slotID})
		if ok {
			callbacks.FinishAll()
		}
	}
}

// inventory contents(for chest...) [NOT TEST]
func (r *Resources) handleInventorySlot(p *packet.InventorySlot) {
	windowID, slotID := WindowID(p.WindowID), SlotID(p.Slot)
	r.inventory.setItemStack(windowID, slotID, &p.NewItem)
	callbacks, ok := r.inventoryCallback.LoadAndDelete(SlotLocation{WindowID: windowID, SlotID: slotID})
	if ok {
		callbacks.FinishAll()
	}
}

// item stack request
func (r *Resources) handleItemStackResponse(p *packet.ItemStackResponse) {
	for _, response := range p.Responses {
		requestID := ItemStackRequestID(response.RequestID)

		callback, ok := r.itemStackCallback.LoadAndDelete(requestID)
		if !ok {
			panic(fmt.Sprintf("handleItemStackResponse: Item stack request with id %d set no callback", response.RequestID))
		}
		containerIDToWindowID, ok := r.itemStackMapping.LoadAndDelete(requestID)
		if !ok {
			panic(fmt.Sprintf("handleItemStackResponse: Item stack request with id %d set no container ID to Window ID mapping", response.RequestID))
		}
		itemUpdater, _ := r.itemStackUpdater.LoadAndDelete(requestID)

		if response.Status != protocol.ItemStackResponseStatusOK {
			callback()
			continue
		}

		for _, containerInfo := range response.ContainerInfo {
			windowID, existed := containerIDToWindowID[ContainerID(containerInfo.ContainerID)]
			if !existed {
				panic(
					fmt.Sprintf(
						"handleItemStackResponse: ContainerID %d not existed in underlying container ID to window ID mapping %#v (request id = %d)",
						containerInfo.ContainerID, containerIDToWindowID, response.RequestID,
					),
				)
			}

			for _, slotInfo := range containerInfo.SlotInfo {
				slotID := SlotID(slotInfo.Slot)

				item, inventoryExisted := r.inventory.GetItemStack(windowID, slotID)
				if !inventoryExisted {
					panic(
						fmt.Sprintf("handleItemStackResponse: Inventory whose window ID is %d is not existed (request id = %d)",
							windowID, response.RequestID,
						),
					)
				}

				UpdateItem(
					item,
					SlotLocation{WindowID: windowID, SlotID: slotID},
					slotInfo, itemUpdater,
				)
				r.inventory.setItemStack(windowID, slotID, item)
			}
		}

		callback()
	}
}

// 根据收到的数据包更新客户端的资源数据
func (r *Resources) handlePacket(pk *packet.Packet) {
	switch p := (*pk).(type) {
	case *packet.CommandOutput:
		r.handleCommandOutput(p)
	case *packet.PyRpc:
		r.handlePyRpc(p)
	case *packet.InventoryContent:
		r.handleInventoryContent(p)
	case *packet.InventoryTransaction:
		r.handleInventoryTransaction(p)
	case *packet.InventorySlot:
		r.handleInventorySlot(p)
	case *packet.ItemStackResponse:
		r.handleItemStackResponse(p)
		// case *packet.ContainerOpen:
		// 	if !r.Container.GetOccupyStates() {
		// 		panic("handlePacket: Attempt to send packet.ContainerOpen without using ResourcesControlCenter")
		// 	}
		// 	r.Container.write_container_closing_data(nil)
		// 	r.Container.write_container_opening_data(p)
		// 	r.Inventory.create_new_inventory(uint32(p.WindowID))
		// 	r.Container.respond_to_container_operation()
		// 	// when a container is opened
		// case *packet.ContainerClose:
		// 	if p.WindowID != 0 && p.WindowID != 119 && p.WindowID != 120 && p.WindowID != 124 {
		// 		err := r.Inventory.delete_inventory(uint32(p.WindowID))
		// 		if err != nil {
		// 			panic(fmt.Sprintf("handlePacket: Try to removed an inventory which not existed; p.WindowID = %v", p.WindowID))
		// 		}
		// 	}
		// 	if !p.ServerSide && !r.Container.GetOccupyStates() {
		// 		panic("handlePacket: Attempt to send packet.ContainerClose without using ResourcesControlCenter")
		// 	}
		// 	r.Container.write_container_opening_data(nil)
		// 	r.Container.write_container_closing_data(p)
		// 	r.Container.respond_to_container_operation()
		// 	// when a container has been closed
		// case *packet.StructureTemplateDataResponse:
		// 	if !r.Structure.GetOccupyStates() {
		// 		panic("handlePacket: Attempt to send packet.StructureTemplateDataRequest without using ResourcesControlCenter")
		// 	}
		// 	r.Structure.writeResponse(*p)
		// 	// used to request mcstructure data
	}
	// // process packet
	// r.Listener.distribute_packet(*pk)
	// // distribute packet(for packet listener)
}
