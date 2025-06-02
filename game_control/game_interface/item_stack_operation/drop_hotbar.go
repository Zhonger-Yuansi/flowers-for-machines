package item_stack_operation

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// DropHotbar 指示物品丢弃操作。
// 不同于 Drop，此操作是针对快捷栏的
type DropHotbar struct {
	Default
	SlotID resources_control.SlotID
	Count  uint8
}

func (DropHotbar) ID() uint8 {
	return IDItemStackOperationDropHotbar
}

func (DropHotbar) CanInline() bool {
	return true
}

func (d DropHotbar) Make(runtimeData MakingRuntime) []protocol.StackRequestAction {
	data := runtimeData.(DropHotbarRuntime)
	return []protocol.StackRequestAction{
		&protocol.DropStackRequestAction{
			Count: d.Count,
			Source: protocol.StackRequestSlotInfo{
				ContainerID:    protocol.ContainerHotBar,
				Slot:           byte(d.SlotID),
				StackNetworkID: data.DropSrcStackNetworkID,
			},
			Randomly: data.Randomly,
		},
	}
}
