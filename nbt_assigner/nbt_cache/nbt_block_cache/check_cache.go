package nbt_block_cache

import (
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
)

// CheckCache 检索整个缓存命中系统，查询 hashNumber 是否存在。
// 返回的 structure 指示命中的结果；
// 返回的 isSetHashHit 指示命中的缓存是否是集合哈希校验和
func (n *NBTBlockCache) CheckCache(hashNumber nbt_hash.CompletelyHashNumber) (
	structure StructureNBTBlock,
	hit bool,
	isSetHashHit bool,
) {
	structure, ok := n.cachedNBTBlock[hashNumber.HashNumber]
	if ok {
		return structure, true, false
	}

	if hashNumber.SetHashNumber == nbt_hash.SetHashNumberNotExist {
		return StructureNBTBlock{}, false, false
	}

	for _, value := range n.cachedNBTBlock {
		if hashNumber.SetHashNumber == value.HashNumber.SetHashNumber {
			return value, true, true
		}
	}

	return StructureNBTBlock{}, false, false
}
