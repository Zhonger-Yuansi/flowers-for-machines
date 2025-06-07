package item_cache

import "github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"

// ConsumeCache 使 slotID 所指示的物品被消耗，
// 这使底层认为 slotID 所指示的槽位可以被重用
func (i *ItemCache) ConsumeCache(slotID resources_control.SlotID) {
	newOne := make([]ItemCacheInfo, 0)

	for _, value := range i.firstCache {
		if value.SlotID == slotID {
			continue
		}
		newOne = append(newOne, value)
	}

	i.firstCache = newOne
	i.console.SetInventorySlot(slotID, false)
}

// ClearnSecondCache 清除 index 所指示方块的二级缓存。
// 它不会改变底层操作台中相应方块的数据
func (i *ItemCache) ClearnSecondCache(index int) {
	i.secondCache[index] = nil
}

// CleanInventory 将背包中所有物品都标记为空气，
// 这使得整个背包的所有物品栏都将可用被重新使用
func (i *ItemCache) CleanInventory() {
	i.firstCache = nil
	i.console.CleanInventory()
}
