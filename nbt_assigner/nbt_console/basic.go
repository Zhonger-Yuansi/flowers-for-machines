package nbt_console

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// API 返回操作台的底层游戏交互接口
func (c *Console) API() *game_interface.GameInterface {
	return c.api
}

// Center 返回操作台中心处的方块坐标
func (c Console) Center() protocol.BlockPos {
	return c.center
}

// Center 返回机器人当前的坐标。不保证是最准确的，
// 因为可能机器人可能会由于受到方块挤压而发生了一定的偏移
func (c Console) Position() protocol.BlockPos {
	return c.position
}

// HotbarSlotID 返回机器人当前所手持物品的快捷栏槽位索引
func (c Console) HotbarSlotID() resources_control.SlotID {
	return c.currentHotBar
}
