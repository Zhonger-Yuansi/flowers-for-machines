package item_cache

import "github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"

// consumeCache ..
func (i *ItemCache) consumeCache(slotID resources_control.SlotID) {
	newOne := make([]ItemCacheInfo, 0)
	for _, value := range i.firstCache {
		if value.SlotID == slotID {
			continue
		}
		newOne = append(newOne, value)
	}
	i.firstCache = newOne
}

// ConsumeCache 使 slotID 所指示的物品被消耗，
// 这使底层认为 slotID 所指示的槽位可以被重用
func (i *ItemCache) ConsumeCache(slotID resources_control.SlotID) {
	i.consumeCache(slotID)
	i.console.UseInventorySlot(i.uniqueID, slotID, false)
}
