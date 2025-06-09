package item_cache

import (
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
)

// CheckCache 检索整个缓存命中系统，查询 hashNumber 是否存在。
// 返回的 isSetHashHit 指示命中的缓存是否是集合哈希校验和
func (i *ItemCache) CheckCache(hashNumber nbt_hash.CompletelyHashNumber) (hit bool, isSetHashHit bool) {
	_, ok := i.thirdCache[hashNumber.HashNumber]
	if ok {
		return true, false
	}

	if hashNumber.SetHashNumber == nbt_hash.SetHashNumberNotExist {
		return false, false
	}

	for _, value := range i.allStructure {
		for _, val := range value.Items {
			if val.ItemInfo.Hash.SetHashNumber == hashNumber.SetHashNumber {
				return true, true
			}
		}
	}

	return false, false
}
