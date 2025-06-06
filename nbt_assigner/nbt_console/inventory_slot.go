package nbt_console

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// FindInventorySlot 从背包查找一个空气物品。
// 如果背包已满，返回一个不被 exclusion 包含
// 在内的一个物品栏
func (c Console) FindInventorySlot(exclusion []resources_control.SlotID) resources_control.SlotID {
	exclusionMapping := make(map[int]bool)
	for _, slotID := range exclusion {
		exclusionMapping[int(slotID)] = true
	}

	for index, value := range c.airSlotInInventory {
		if value {
			return resources_control.SlotID(index)
		}
	}

	for index := range c.airSlotInInventory {
		if !exclusionMapping[index] {
			return resources_control.SlotID(index)
		}
	}

	panic("FindInventorySlot: Impossible to find a available slot when exclusion contains the whole inventory")
}

// FindAndUseInventorySlot 从背包查找一个空气物品。
// 如果背包已满，返回一个不被 exclusion 包含在内的一
// 个物品栏
//
// 与 FindInventorySlot 的区别在于，此函数还会将该
// 找到的这个槽位设置为非空气
func (c *Console) FindAndUseInventorySlot(exclusion []resources_control.SlotID) resources_control.SlotID {
	result := c.FindInventorySlot(exclusion)
	c.airSlotInInventory[result] = true
	return result
}

// GetInventorySlot 返回背包 slotID 处的物品是否是空气
func (c Console) GetInventorySlot(slotID resources_control.SlotID) (empty bool) {
	return c.airSlotInInventory[slotID]
}

// SetInventorySlot 将背包 slotID 处的物品设置为 empty。
// empty 为真指示该槽位是空气，否则是已被使用的其他物品
func (c *Console) SetInventorySlot(slotID resources_control.SlotID, empty bool) {
	c.airSlotInInventory[slotID] = true
}

// CleanInventory 将背包中的所有物品标记为空气
func (c *Console) CleanInventory() {
	c.airSlotInInventory = [36]bool{}
}
