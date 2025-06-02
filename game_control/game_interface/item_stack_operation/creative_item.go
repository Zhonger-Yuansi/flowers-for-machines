package item_stack_operation

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// CreativeItem 指示创造物品获取操作
type CreativeItem struct {
	Default
	CreativeItemNetworkID uint32
	SlotID                resources_control.SlotID
	Count                 uint8
}

func (CreativeItem) ID() uint8 {
	return IDItemStackOperationCreativeItem
}

func (CreativeItem) CanInline() bool {
	return false
}

func (d CreativeItem) Make(runtimeData MakingRuntime) []protocol.StackRequestAction {
	data := runtimeData.(CreativeItemRuntime)

	move := protocol.PlaceStackRequestAction{}
	move.Count = d.Count
	move.Source = protocol.StackRequestSlotInfo{
		ContainerID:    protocol.ContainerCreatedOutput,
		Slot:           0x32,
		StackNetworkID: data.RequestID,
	}
	move.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    protocol.ContainerCombinedHotBarAndInventory,
		Slot:           byte(d.SlotID),
		StackNetworkID: 0,
	}

	return []protocol.StackRequestAction{
		&protocol.CraftCreativeStackRequestAction{CreativeItemNetworkID: d.CreativeItemNetworkID},
		&move,
	}
}
