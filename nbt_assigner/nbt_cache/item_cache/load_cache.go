package item_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// LoadCache 加载哈希校验和为 hashNumber 的物品到物品栏。
//
// 提供的 exclusion 确保它其中所示的背包物品栏不会被意外
// 使用，应当确保这些物品栏都装有物品
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
		hit, isSetHashHit = false, false
	}
	return
}
