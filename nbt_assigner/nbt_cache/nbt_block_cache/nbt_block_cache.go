package nbt_block_cache

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/google/uuid"
)

// NBTBlockCache 是基于操作台实现的 NBT 方块缓存命中系统
type NBTBlockCache struct {
	// uniqueID 是当前缓存命中系统的唯一标识
	uniqueID string
	// console 是机器人使用的操作台
	console *nbt_console.Console
	// cachedNBTBlock 记载了已缓存的所有 NBT 方块
	cachedNBTBlock map[uint64]StructureNBTBlock
}

// NewNBTBlockCache 基于操作台 console 创建并返回一个新的 NBT 方块缓存命中系统
func NewNBTBlockCache(console *nbt_console.Console) *NBTBlockCache {
	return &NBTBlockCache{
		uniqueID:       uuid.NewString(),
		console:        console,
		cachedNBTBlock: make(map[uint64]StructureNBTBlock),
	}
}
