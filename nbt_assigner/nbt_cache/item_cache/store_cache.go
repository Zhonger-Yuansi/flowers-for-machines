package item_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/google/uuid"
)

// StoreCache 保存操作台中心方块所指示容器中的全部物品到当前缓存命中系统。
// items 是该容器中所有物品的数据，container 则指示应当如何打开这个容器
func (i *ItemCache) StoreCache(
	items []ItemCacheInfo,
	container block_helper.ContainerBlockOpenInfo,
) error {
	if len(items) == 0 {
		return nil
	}

	block := i.console.BlockByIndex(nbt_console.ConsoleIndexCenterBlock)
	if _, ok := (*block).(block_helper.ContainerBlockHelper); !ok {
		return fmt.Errorf(
			"StoreCache: Center block who at %#v is not a container; block = %#v",
			i.console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock), block,
		)
	}

	allHit := true
	for _, item := range items {
		if _, ok := i.thirdCache[item.Hash.HashNumber]; !ok {
			allHit = false
			break
		}
	}
	if allHit {
		return nil
	}

	uniqueID, err := i.console.API().StructureBackup().BackupStructure(i.console.Center())
	if err != nil {
		return fmt.Errorf("StoreCache: %v", err)
	}

	for _, item := range items {
		if _, ok := i.thirdCache[item.Hash.HashNumber]; ok {
			continue
		}
		i.thirdCache[item.Hash.HashNumber] = StructureItemCache{
			UniqueID: uniqueID,
			CompletelyInfo: CompletelyItemInfo{
				ContainerInfo: container,
				ItemInfo:      item,
			},
		}
	}

	structureItems := StructureItems{
		ContainerInfo: container,
	}
	for _, item := range items {
		structureItems.Items = append(structureItems.Items, CompletelyItemInfo{
			ContainerInfo: container,
			ItemInfo:      item,
		})
	}
	i.allStructure[uniqueID] = structureItems

	return nil
}

// CleanCache 清除三级缓存中的所有内容。
// 它从设计上确保不会返回错误
func (i *ItemCache) CleanThirdCache() {
	api := i.console.API().StructureBackup()

	for _, value := range i.thirdCache {
		_ = api.DeleteStructure(value.UniqueID)
	}

	i.allStructure = make(map[uuid.UUID]StructureItems)
	i.thirdCache = make(map[int64]StructureItemCache)
}
