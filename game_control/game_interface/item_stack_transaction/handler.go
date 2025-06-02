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
	virtualInventories *virtualInventories
	responseMapping    *responseMapping
}

// newItemStackOperationHandler ..
func newItemStackOperationHandler(
	api *resources_control.ContainerManager,
	virtualInventories *virtualInventories,
	responseMapping *responseMapping,
) *itemStackOperationHandler {
	return &itemStackOperationHandler{
		api:                api,
		virtualInventories: virtualInventories,
		responseMapping:    responseMapping,
	}
}

// handleMove ..
func (i *itemStackOperationHandler) handleMove(
	op item_stack_operation.Move,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
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

	runtimeData := item_stack_operation.DropRuntime{
		DropSrcContainerID:    byte(srcCID),
		DropSrcStackNetworkID: srcRID,
		Randomly:              false,
	}
	return op.Make(runtimeData), nil
}

// handleDropHotbar ..
func (i *itemStackOperationHandler) handleDropHotbar(
	op item_stack_operation.DropHotbar,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	slot := resources_control.SlotLocation{
		WindowID: protocol.WindowIDInventory,
		SlotID:   op.SlotID,
	}

	srcRID, err := i.virtualInventories.loadAndSetStackNetworkID(slot, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleDropHotbar: %v", err)
	}
	i.responseMapping.bind(protocol.WindowIDInventory, protocol.ContainerHotBar)

	runtimeData := item_stack_operation.DropHotbarRuntime{
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
	i.responseMapping.bind(protocol.WindowIDInventory, protocol.ContainerCombinedHotBarAndInventory)
	return op.Make(
		item_stack_operation.CreativeItemRuntime{
			RequestID: int32(requestID),
		},
	), nil
}

// handleRenaming ..
func (i *itemStackOperationHandler) handleRenaming(
	op item_stack_operation.Renaming,
	requestID resources_control.ItemStackRequestID,
) (result []protocol.StackRequestAction, err error) {
	itemCount, err := i.virtualInventories.loadItemCount(op.Path)
	if err != nil {
		return nil, fmt.Errorf("handleRenaming: %v", err)
	}

	srcRID, err := i.virtualInventories.loadAndSetStackNetworkID(op.Path, requestID)
	if err != nil {
		return nil, fmt.Errorf("handleRenaming: %v", err)
	}

	srcCID, found := slotLocationToContainerID(i.api, op.Path)
	if !found {
		return nil, fmt.Errorf("handleRenaming: Can not find the container ID of given item whose at %#v", op.Path)
	}

	containerData, existed := i.api.ContainerData()
	if !existed {
		return nil, fmt.Errorf("handleRenaming: Anvil is not opened")
	}

	i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerAnvilInput)
	i.responseMapping.bind(op.Path.WindowID, srcCID)

	runtimeData := item_stack_operation.RenamingRuntime{
		RequestID:      int32(requestID),
		ItemCount:      uint8(itemCount),
		ContainerID:    byte(srcCID),
		StackNetworkID: srcRID,
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

	containerData, existed := i.api.ContainerData()
	if !existed {
		return nil, fmt.Errorf("handleLooming: Loom is not opened")
	}

	if op.UsePattern {
		rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.PatternPath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		cid, found := slotLocationToContainerID(i.api, op.PatternPath)
		if !found {
			return nil, fmt.Errorf("handleLooming: Can not find the container ID of given item whose at %#v", op.PatternPath)
		}

		i.responseMapping.bind(op.PatternPath.WindowID, cid)
		i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerLoomMaterial)

		runtimeData.MovePatternSrcContainerID = byte(cid)
		runtimeData.MovePatternSrcStackNetworkID = rid
	}

	// Banner
	{
		rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.BannerPath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		cid, found := slotLocationToContainerID(i.api, op.BannerPath)
		if !found {
			return nil, fmt.Errorf("handleLooming: Can not find the container ID of given item whose at %#v", op.BannerPath)
		}

		i.responseMapping.bind(op.BannerPath.WindowID, cid)
		i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerLoomInput)

		runtimeData.MoveBannerSrcContainerID = byte(cid)
		runtimeData.MoveBannerSrcStackNetworkID = rid
	}

	// Dye
	{
		rid, err := i.virtualInventories.loadAndSetStackNetworkID(op.DyePath, requestID)
		if err != nil {
			return nil, fmt.Errorf("handleLooming: %v", err)
		}

		cid, found := slotLocationToContainerID(i.api, op.DyePath)
		if !found {
			return nil, fmt.Errorf("handleLooming: Can not find the container ID of given item whose at %#v", op.DyePath)
		}

		i.responseMapping.bind(op.DyePath.WindowID, cid)
		i.responseMapping.bind(resources_control.WindowID(containerData.WindowID), protocol.ContainerLoomDye)

		runtimeData.MoveDyeSrcContainerID = byte(cid)
		runtimeData.MoveDyeSrcStackNetworkID = rid
	}

	return op.Make(runtimeData), nil
}
