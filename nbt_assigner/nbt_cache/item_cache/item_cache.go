package item_cache

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/google/uuid"
)

// ItemCache 是基于操作台实现的物品缓存命中系统
type ItemCache struct {
	// uniqueID 是当前缓存命中系统的唯一标识
	uniqueID string
	// console 是机器人使用的操作台
	console *nbt_console.Console

	// allStructure 记载每个结构中物品的分布。
	// 应确保这里的每个结构都保存的是一个不相
	// 连的容器
	allStructure map[uuid.UUID]StructureItems
	// thirdCache 是缓存命中系统的第三级缓存，
	// 这意味着物品位于结构命令保存的结构中
	thirdCache map[uint64]StructureItemCache

	// secondCache 是缓存命中系统的第二级缓存，
	// 这意味着物品已经出现在操作台附近的帮助方块中
	secondCache [9][]CompletelyItemInfo
	// firstCache 是缓存命中系统的第一级缓存，
	// 这意味着物品已经被载入到机器人的背包中
	firstCache []ItemCacheInfo
}

// NewItemCache 基于操作台 console 创建并返回一个新的物品缓存命中系统
func NewItemCache(console *nbt_console.Console) *ItemCache {
	itemCache := &ItemCache{
		uniqueID:     uuid.NewString(),
		console:      console,
		allStructure: make(map[uuid.UUID]StructureItems),
		thirdCache:   make(map[uint64]StructureItemCache),
		secondCache:  [9][]CompletelyItemInfo{},
		firstCache:   nil,
	}

	console.SetHelperUseCallback(itemCache.useBlocksCallback)
	console.SetSlotUseCallback(itemCache.useSlotCallback)

	return itemCache
}
