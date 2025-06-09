package nbt_block_cache

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
)

// CheckCache 检索整个缓存命中系统，查询 hashNumber 是否存在。
// 返回的 offset 指示 NBT 方块的偏移，例如床尾相对于床头的偏移；
// 返回的 isSetHashHit 指示命中的缓存是否是集合哈希校验和；
func (n *NBTBlockCache) CheckCache(hashNumber nbt_hash.CompletelyHashNumber) (offset protocol.BlockPos, hit bool, isSetHashHit bool) {
	structure, ok := n.cachedNBTBlock[hashNumber.HashNumber]
	if ok {
		return structure.Offset, true, false
	}

	if hashNumber.SetHashNumber == nbt_hash.SetHashNumberNotExist {
		return protocol.BlockPos{}, false, false
	}

	for _, value := range n.cachedNBTBlock {
		if hashNumber.SetHashNumber == value.HashNumber.SetHashNumber {
			return value.Offset, true, true
		}
	}

	return protocol.BlockPos{}, false, false
}
