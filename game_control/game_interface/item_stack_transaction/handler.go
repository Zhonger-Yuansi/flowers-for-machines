package item_stack_transaction

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface/item_stack_operation"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// itemStackOperationHandler ..
type itemStackOperationHandler struct {
	api                *resources_control.ContainerManager
	constantPacket     *resources_control.ConstantPacket
	virtualInventories *virtualInventories
	responseMapping    *responseMapping
}

// newItemStackOperationHandler ..
func newItemStackOperationHandler(
	api *resources_control.ContainerManager,
	constantPacket *resources_control.ConstantPacket,
	virtualInventories *virtualInventories,
	responseMapping *responseMapping,
) *itemStackOperationHandler {
	return &itemStackOperationHandler{
		api:                api,
		constantPacket:     constantPacket,
		virtualInventories: virtualInventories,
		responseMapping:    responseMapping,
	}
}

// handleMove ..
func (i *itemStackOperationHandler) handleMove(
	op item_stack_operation.Move,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	if op.Source == op.Destination {
		return nil, fmt.Errorf("handleMove: Source is equal to Destination")
	}

	srcRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Source, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleMove: %v", err)
	}
	dstRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Destination, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleMove: %v", err)
	}

	srcCID, found := slotLocationToContainerID(i.api, op.Source)
	if !found {
		return nil, fmt.Errorf("handleMove: Can not find the container ID of given item whose at %#v", op.Source)
	}
	dstCID, found := slotLocationToContainerID(i.api, op.Destination)
	if !found {
		return nil, fmt.Errorf("handleMove: Can not find the container ID of given item whose at %#v", op.Destination)
	}
	i.responseMapping.bind(op.Source.WindowID, srcCID)
	i.responseMapping.bind(op.Destination.WindowID, dstCID)

	_, err = i.virtualInventories.loadAndAddItemCount(op.Source, -int8(op.Count), false)
	if err != nil {
		return nil, fmt.Errorf("handleMove: %v", err)
	}
	_, err = i.virtualInventories.loadAndAddItemCount(op.Destination, int8(op.Count), false)
	if err != nil {
		return nil, fmt.Errorf("handleMove: %v", err)
	}

	runtimeData := item_stack_operation.MoveRuntime{
		MoveSrcContainerID:    byte(srcCID),
		MoveSrcStackNetworkID: srcRID,
		MoveDstContainerID:    byte(dstCID),
		MoveDstStackNetworkID: dstRID,
	}
	return op.Make(runtimeData), nil
}

// handleSwap ..
func (i *itemStackOperationHandler) handleSwap(
	op item_stack_operation.Swap,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	if op.Source == op.Destination {
		return nil, fmt.Errorf("handleSwap: Source is equal to Destination")
	}

	srcRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Source, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleSwap: %v", err)
	}
	dstRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Destination, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleSwap: %v", err)
	}

	srcCID, found := slotLocationToContainerID(i.api, op.Source)
	if !found {
		return nil, fmt.Errorf("handleSwap: Can not find the container ID of given item whose at %#v", op.Source)
	}
	dstCID, found := slotLocationToContainerID(i.api, op.Destination)
	if !found {
		return nil, fmt.Errorf("handleSwap: Can not find the container ID of given item whose at %#v", op.Destination)
	}
	i.responseMapping.bind(op.Source.WindowID, srcCID)
	i.responseMapping.bind(op.Destination.WindowID, dstCID)

	srcCount, err := i.virtualInventories.loadAndAddItemCount(op.Source, 0, true)
	if err != nil {
		return nil, fmt.Errorf("handleSwap: %v", err)
	}
	dstCount, err := i.virtualInventories.loadAndAddItemCount(op.Destination, 0, true)
	if err != nil {
		return nil, fmt.Errorf("handleSwap: %v", err)
	}
	_, err = i.virtualInventories.loadAndSetItemCount(op.Source, dstCount)
	if err != nil {
		return nil, fmt.Errorf("handleSwap: %v", err)
	}
	_, err = i.virtualInventories.loadAndSetItemCount(op.Destination, srcCount)
	if err != nil {
		return nil, fmt.Errorf("handleSwap: %v", err)
	}

	runtimeData := item_stack_operation.SwapRuntime{
		SwapSrcContainerID:    byte(srcCID),
		SwapSrcStackNetworkID: srcRID,
		SwapDstContainerID:    byte(dstCID),
		SwapDstStackNetworkID: dstRID,
	}
	return op.Make(runtimeData), nil
}

// handleDrop ..
func (i *itemStackOperationHandler) handleDrop(
	op item_stack_operation.Drop,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	srcRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Path, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleDrop: %v", err)
	}

	srcCID, found := slotLocationToContainerID(i.api, op.Path)
	if !found {
		return nil, fmt.Errorf("handleDrop: Can not find the container ID of given item whose at %#v", op.Path)
	}
	i.responseMapping.bind(op.Path.WindowID, srcCID)

	_, err = i.virtualInventories.loadAndAddItemCount(op.Path, -int8(op.Count), false)
	if err != nil {
		return nil, fmt.Errorf("handleDrop: %v", err)
	}

	runtimeData := item_stack_operation.DropRuntime{
		DropSrcContainerID:    byte(srcCID),
		DropSrcStackNetworkID: srcRID,
		Randomly:              false,
	}
	return op.Make(runtimeData), nil
}

// handleCreativeItem ..
func (i *itemStackOperationHandler) handleCreativeItem(
	op item_stack_operation.CreativeItem,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	var creativeItemNetworkID uint32

	rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.Path, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleCreativeItem: %v", err)
	}

	cid, found := slotLocationToContainerID(i.api, op.Path)
	if !found {
		return nil, fmt.Errorf("handleCreativeItem: Can not find the container ID of given item whose at %#v", op.Path)
	}
	i.responseMapping.bind(op.Path.WindowID, cid)

	_, err = i.virtualInventories.loadAndAddItemCount(op.Path, int8(op.Count), false)
	if err != nil {
		return nil, fmt.Errorf("handleCreativeItem: %v", err)
	}

	if op.UseCreativeItemNetworkID {
		creativeItemNetworkID = op.CreativeItemNetworkID
	}
	if op.UseNetworkID {
		creativeItemNetworkID = i.constantPacket.CreativeItemByNI(op.NetworkID).CreativeItemNetworkID
	}

	return op.Make(
		item_stack_operation.CreativeItemRuntime{
			RequestID:             int32(requestID),
			DstContainerID:        byte(cid),
			DstItemStackID:        rid,
			CreativeItemNetworkID: creativeItemNetworkID,
		},
	), nil
}

// handleRenaming ..
func (i *itemStackOperationHandler) handleRenaming(
	op item_stack_operation.Renaming,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	containerData, _, existed := i.api.ContainerData()
	if !existed {
		return nil, fmt.Errorf("handleRenaming: Anvil is not opened")
	}

	srcRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Path, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleRenaming: %v", err)
	}
	anvilRID, err := i.virtualInventories.loadAndSetStackNetworkID(
		resources_control.SlotLocation{
			WindowID: resources_control.WindowID(containerData.WindowID),
			SlotID:   1,
		},
		requestID,
	)
	if err != nil {
		return nil, fmt.Errorf("handleRenaming: %v", err)
	}

	srcCID, found := slotLocationToContainerID(i.api, op.Path)
	if !found {
		return nil, fmt.Errorf("handleRenaming: Can not find the container ID of given item whose at %#v", op.Path)
	}
	i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerAnvilInput)
	i.responseMapping.bind(op.Path.WindowID, srcCID)

	srcCount, err := i.virtualInventories.loadAndAddItemCount(op.Path, 0, true)
	if err != nil {
		return nil, fmt.Errorf("handleRenaming: %v", err)
	}

	runtimeData := item_stack_operation.RenamingRuntime{
		RequestID:               int32(requestID),
		ItemCount:               srcCount,
		SrcContainerID:          byte(srcCID),
		SrcStackNetworkID:       srcRID,
		AnvilSlotStackNetworkID: anvilRID,
	}
	return op.Make(runtimeData), nil
}

// handleLooming ..
func (i *itemStackOperationHandler) handleLooming(
	op item_stack_operation.Looming,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	runtimeData := item_stack_operation.LoomingRuntime{
		RequestID: int32(requestID),
	}

	if op.BannerPath == op.DyePath {
		return nil, fmt.Errorf("handleSwap: BannerPath is equal to DyePath")
	}
	if op.UsePattern {
		if op.PatternPath == op.BannerPath {
			return nil, fmt.Errorf("handleSwap: PatternPath is equal to BannerPath")
		}
		if op.PatternPath == op.DyePath {
			return nil, fmt.Errorf("handleSwap: PatternPath is equal to DyePath")
		}
	}

	containerData, _, existed := i.api.ContainerData()
	if !existed {
		return nil, fmt.Errorf("handleLooming: Loom is not opened")
	}

	if op.UsePattern {
		loomSlot := resources_control.SlotLocation{
			WindowID: resources_control.WindowID(containerData.WindowID),
			SlotID:   11,
		}

		rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.PatternPath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}
		loomRID, err := i.virtualInventories.loadAndSetStackNetworkID(loomSlot, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		cid, found := slotLocationToContainerID(i.api, op.PatternPath)
		if !found {
			return nil, fmt.Errorf("handleLooming: Can not find the container ID of given item whose at %#v", op.PatternPath)
		}
		i.responseMapping.bind(op.PatternPath.WindowID, cid)
		i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerLoomMaterial)

		_, err = i.virtualInventories.loadAndAddItemCount(op.PatternPath, 0, true)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		runtimeData.LoomPatternStackNetworkID = loomRID
		runtimeData.MovePatternSrcContainerID = byte(cid)
		runtimeData.MovePatternSrcStackNetworkID = rid
	}

	// Banner
	{
		loomSlot := resources_control.SlotLocation{
			WindowID: resources_control.WindowID(containerData.WindowID),
			SlotID:   9,
		}

		rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.BannerPath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}
		loomRID, err := i.virtualInventories.loadAndSetStackNetworkID(loomSlot, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		cid, found := slotLocationToContainerID(i.api, op.BannerPath)
		if !found {
			return nil, fmt.Errorf("handleLooming: Can not find the container ID of given item whose at %#v", op.BannerPath)
		}
		i.responseMapping.bind(op.BannerPath.WindowID, cid)
		i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerLoomInput)

		_, err = i.virtualInventories.loadAndAddItemCount(op.BannerPath, 0, true)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		runtimeData.LoomBannerStackNetworkID = loomRID
		runtimeData.MoveBannerSrcContainerID = byte(cid)
		runtimeData.MoveBannerSrcStackNetworkID = rid
	}

	// Dye
	{
		loomSlot := resources_control.SlotLocation{
			WindowID: resources_control.WindowID(containerData.WindowID),
			SlotID:   10,
		}

		rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.DyePath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}
		loomRID, err := i.virtualInventories.loadAndSetStackNetworkID(loomSlot, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		cid, found := slotLocationToContainerID(i.api, op.DyePath)
		if !found {
			return nil, fmt.Errorf("handleLooming: Can not find the container ID of given item whose at %#v", op.DyePath)
		}
		i.responseMapping.bind(op.DyePath.WindowID, cid)
		i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerLoomDye)

		_, err = i.virtualInventories.loadAndAddItemCount(op.DyePath, -1, false)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		runtimeData.LoomDyeStackNetworkID = loomRID
		runtimeData.MoveDyeSrcContainerID = byte(cid)
		runtimeData.MoveDyeSrcStackNetworkID = rid
	}

	return op.Make(runtimeData), nil
}

// handleCrafting ..
func (i *itemStackOperationHandler) handleCrafting(
	op item_stack_operation.Crafting,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	runtimeData := item_stack_operation.CraftingRuntime{
		RequestID: int32(requestID),
	}

	for slotID, count := range i.virtualInventories.allItemCount(protocol.WindowIDCrafting) {
		if count == 0 {
			continue
		}

		location := resources_control.SlotLocation{
			WindowID: protocol.WindowIDCrafting,
			SlotID:   slotID,
		}

		rid, err := i.virtualInventories.loadAndSetStackNetworkID(location, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleCrafting: %v", err)
		}

		err = i.virtualInventories.setItemCount(location, 0)
		if err != nil {
			return nil, fmt.Errorf("handleCrafting: %v", err)
		}

		runtimeData.Consumes = append(runtimeData.Consumes, item_stack_operation.CraftingConsume{
			Slot:           location.SlotID,
			StackNetworkID: rid,
			Count:          count,
		})
	}

	{
		resultPath := resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   op.ResultSlotID,
		}

		rid, err := i.virtualInventories.loadAndSetStackNetworkID(resultPath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleCrafting: %v", err)
		}

		i.virtualInventories.addItemCount(resultPath, int8(op.ResultCount), false)
		runtimeData.ResultStackNetworkID = rid
	}

	i.responseMapping.bind(protocol.WindowIDInventory, protocol.ContainerCombinedHotBarAndInventory)
	i.responseMapping.bind(protocol.WindowIDCrafting, protocol.ContainerCraftingInput)

	return op.Make(runtimeData), nil
}
