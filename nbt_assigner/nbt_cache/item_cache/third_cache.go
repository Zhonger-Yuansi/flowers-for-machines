package item_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
)

// loadThirdCacheToSecond 从三级缓存 (已保存的结构) 中加载 hashNumber 所示的物品到二级缓存上。
// isSetHashHit 指示命中的缓存是否是不完整的，而只是命中了集合哈希校验和
func (i *ItemCache) loadThirdCacheToSecond(hashNumber ItemHashNumber) (hit bool, isSetHashHit bool, err error) {
	var structure StructureItemCache
	var pos protocol.BlockPos
	api := i.console.API().StructureBackup()

	// Try completely hash number
	structure, hit = i.thirdCache[hashNumber.HashNumber]
	if !hit && hashNumber.SetHashNumber != SetHashNumberNotExist {
		// Try set hash number
		for _, value := range i.thirdCache {
			if value.CompletelyInfo.ItemInfo.Hash.SetHashNumber == hashNumber.SetHashNumber {
				structure, isSetHashHit, hit = value, true, true
				break
			}
		}
	}

	// No cache was hit
	if !hit {
		return false, false, nil
	}

	// Find new space to load cache
	index, _, block := i.console.FindSpaceToPlaceNewContainer(false, true)
	if block != nil {
		index = nbt_console.ConsoleIndexFirstHelperBlock
		block = i.console.BlockByIndex(index)
	}
	pos = i.console.BlockPosByIndex(index)

	// Load cache
	err = api.RevertStructure(structure.UniqueID, pos)
	if err != nil {
		return false, false, fmt.Errorf("loadThirdCacheToSecond: %v", err)
	}

	// Update underlying container data
	container := block_helper.ContainerBlockHelper{
		OpenInfo: structure.CompletelyInfo.ContainerInfo,
		IsEmpty:  false,
	}
	i.console.UseHelperBlock(i.uniqueID, index, container)

	// Update second cache data
	newOne := make([]CompletelyItemInfo, 0)
	for _, value := range i.allStructure[structure.UniqueID].Items {
		newOne = append(newOne, CompletelyItemInfo{
			ContainerInfo: value.ContainerInfo,
			ItemInfo:      value.ItemInfo,
		})
	}
	i.secondCache[index] = newOne

	return true, isSetHashHit, nil
}
