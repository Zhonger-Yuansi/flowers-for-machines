package item_cache

import "github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"

// useBlocksCallback ..
func (i *ItemCache) useBlocksCallback(requester string, index int) {
	if requester == i.uniqueID {
		return
	}
	i.secondCache[index] = nil
}

// useSlotCallback ..
func (i *ItemCache) useSlotCallback(requester string, slotID resources_control.SlotID) {
	if requester == i.uniqueID {
		return
	}
	i.consumeCache(slotID)
}
