package nbt_cache

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache/base_container_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache/item_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache/nbt_block_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
)

// NBTCacheSystem 是基于操作台实现的 NBT 缓存命中系统
type NBTCacheSystem struct {
	i *item_cache.ItemCache
	b *base_container_cache.BaseContainerCache
	n *nbt_block_cache.NBTBlockCache
}

// NewNBTCacheSystem 基于操作台 console 创建并返回一个新的 NBT 缓存命中系统
func NewNBTCacheSystem(console *nbt_console.Console) *NBTCacheSystem {
	return &NBTCacheSystem{
		i: item_cache.NewItemCache(console),
		b: base_container_cache.NewBaseContainerCache(console),
		n: nbt_block_cache.NewNBTBlockCache(console),
	}
}

// ItemCache 返回物品缓存命中系统
func (n *NBTCacheSystem) ItemCache() *item_cache.ItemCache {
	return n.i
}

// BaseContainerCache 返回基容器缓存命中系统
func (n *NBTCacheSystem) BaseContainerCache() *base_container_cache.BaseContainerCache {
	return n.b
}

// NBTBlockCache 返回 NBT 方块缓存命中系统
func (n *NBTCacheSystem) NBTBlockCache() *nbt_block_cache.NBTBlockCache {
	return n.n
}
