package item_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// loadSecondCacheToFirstCache 从二级缓存 (容器)
// 加载校验和为 hashNumber 的物品到一级缓存 (背包)。
//
// 返回的 isSetHashHit 指示命中的缓存是否是不完整的，
// 而只是命中了集合哈希校验和。exclusion 指示排除列表，
// 它保证列表内包含的所有背包物品栏不会被意外使用。
//
// 在设计上应当保证 exclusion 的长度最大为 27，
// 否则这被认为是设计者的非正常使用，并可能伴随程序崩溃
func (i *ItemCache) loadSecondCacheToFirstCache(
	hashNumber ItemHashNumber,
	exclusion []resources_control.SlotID,
) (hit bool, isSetHashHit bool, err error) {
	var hitContainerIndex int
	var hitSliceIndex int
	var hitItem ItemCacheInfo
	var inventorySlot resources_control.SlotID
	api := i.console.API()

	// Check hash is hit or load cache from third cache
	for range 2 {
		// Firstly, check the container we already loaded
		for index, container := range i.secondCache {
			// Only set hash number hit
			if hashNumber.SetHashNumber != SetHashNumberNotExist {
				for idx, item := range container {
					if item.ItemInfo.Hash.SetHashNumber == hashNumber.SetHashNumber {
						hit, isSetHashHit, hitItem = true, true, item.ItemInfo
						hitContainerIndex, hitSliceIndex = index, idx
						break
					}
				}
			}
			// Completely hit
			for idx, item := range container {
				if item.ItemInfo.Hash.HashNumber == hashNumber.HashNumber {
					hit, hitItem = true, item.ItemInfo
					hitContainerIndex, hitSliceIndex = index, idx
					break
				}
			}
			// If hit, then break
			if hit {
				break
			}
		}
		// if hit is false, then the second cache have no item to hit,
		// so here we try to load cache from thrid cache.
		if hit {
			break
		}
		// If we hit and load something from thrid cache to the second cache,
		// then the next loop will hit from the second cache.
		// Due to we loop for two times at most, so if we hit here, and then
		// the loop will break at next time.
		hit, isSetHashHit, err = i.loadThirdCacheToSecond(hashNumber)
		if err != nil {
			return false, false, fmt.Errorf("loadSecondCacheToFirstCache: %v", err)
		}
		if !hit {
			return false, false, nil
		}
	}

	// Open container
	success, err := i.console.OpenContainerByIndex(hitContainerIndex)
	if err != nil {
		return false, false, fmt.Errorf("loadSecondCacheToFirstCache: %v", err)
	}
	if !success {
		return false, false, fmt.Errorf(
			"loadSecondCacheToFirstCache: Failed to open the container which is %#v",
			*i.console.BlockByIndex(hitContainerIndex),
		)
	}

	// Move item / Load cache from second cache
	{
		// Find a possible place to place the cached item
		inventorySlot = i.console.FindInventorySlot(exclusion)
		// If the inventory is full, then we try to grow a new air to place the item
		if !i.console.GetInventorySlot(inventorySlot) {
			err = api.Replaceitem().ReplaceitemInInventory(
				"@s",
				game_interface.ReplacePathInventory,
				game_interface.ReplaceitemInfo{
					Name:     "minecraft:air",
					Count:    1,
					MetaData: 0,
					Slot:     inventorySlot,
				},
				"",
				true,
			)
			if err != nil {
				return false, false, fmt.Errorf("loadSecondCacheToFirstCache: %v", err)
			}
			i.console.SetInventorySlot(inventorySlot, true)
		}
		// Load cache from the helper container block
		success, _, _, err = api.ItemStackOperation().OpenTransaction().
			MoveToInventory(hitItem.SlotID, inventorySlot, 1).
			Commit()
		if err != nil {
			return false, false, fmt.Errorf("loadSecondCacheToFirstCache: %v", err)
		}
		if !success {
			return false, false, fmt.Errorf("loadSecondCacheToFirstCache: %v", err)
		}
		i.console.SetInventorySlot(inventorySlot, false)
	}

	// Close container
	err = api.ContainerOpenAndClose().CloseContainer()
	if err != nil {
		return false, false, fmt.Errorf("loadSecondCacheToFirstCache: %v", err)
	}

	// Update second cache data
	i.secondCache[hitContainerIndex][hitSliceIndex].ItemInfo.Count -= 1
	if i.secondCache[hitContainerIndex][hitSliceIndex].ItemInfo.Count == 0 {
		newOne := make([]CompletelyItemInfo, 0)
		for index, value := range i.secondCache[hitContainerIndex] {
			if index == hitSliceIndex {
				continue
			}
			newOne = append(newOne, value)
		}
		i.secondCache[hitContainerIndex] = newOne
	}

	// Update first cache data
	{
		hitFirstCache := false
		itemCacheInfo := ItemCacheInfo{
			SlotID: inventorySlot,
			Count:  1,
			Hash:   hitItem.Hash,
		}

		for index, value := range i.firstCache {
			if value.SlotID == inventorySlot {
				i.firstCache[index] = itemCacheInfo
				hitFirstCache = true
				break
			}
		}

		if !hitFirstCache {
			i.firstCache = append(i.firstCache, itemCacheInfo)
		}
	}

	return true, isSetHashHit, nil
}
