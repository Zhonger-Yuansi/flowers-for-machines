package item_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// LoadCache 加载哈希校验和为 hashNumber 的物品到物品栏。
// 提供的 exclusion 确保它其中所示的背包物品栏不会被意外使用。
//
// 返回的 slotID 指示加载完成后物品所在的背包槽位索引；
// 返回的 isSetHashHit 指示命中的缓存是否是不完整的，
// 而只是命中了集合哈希校验和。
//
// 在使用 slotID 处的物品后应当立即使用 ConsumeCache 消耗它
func (i *ItemCache) LoadCache(hashNumber ItemHashNumber, exclusion []resources_control.SlotID) (
	slotID resources_control.SlotID,
	hit bool,
	isSetHashHit bool,
	err error,
) {
	for range 2 {
		// Completely hit
		for _, value := range i.firstCache {
			if value.Hash.HashNumber == hashNumber.HashNumber {
				hit, slotID = true, value.SlotID
				return
			}
		}
		// Only set hash number hit
		if hashNumber.SetHashNumber != SetHashNumberNotExist {
			for _, value := range i.firstCache {
				if value.Hash.SetHashNumber == hashNumber.SetHashNumber {
					hit, slotID, isSetHashHit = true, value.SlotID, true
					return
				}
			}
		}
		// Try load from second cache
		hit, isSetHashHit, err = i.loadSecondCacheToFirstCache(hashNumber, exclusion)
		if err != nil {
			return 0, false, false, fmt.Errorf("LoadCache: %v", err)
		}
		if !hit {
			return
		}
	}
	return
}

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
	i.console.SetInventorySlot(slotID, true)
}

// CleanInventory 将背包中所有物品都标记为空气，
// 这使得整个背包的所有物品栏都将可用被重新使用
func (i *ItemCache) CleanInventory() {
	i.firstCache = nil
	i.console.CleanInventory()
}
