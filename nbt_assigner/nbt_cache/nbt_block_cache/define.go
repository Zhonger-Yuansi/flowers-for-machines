package nbt_block_cache

import (
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/google/uuid"
)

// StructureNBTBlock 指示了一个保存在结构中的 NBT 方块
type StructureNBTBlock struct {
	// UniqueID 是这个方块的唯一标识符
	UniqueID uuid.UUID
	// HashNumber 是这个方块的哈希校验和
	HashNumber nbt_hash.CompletelyHashNumber
	// Block 是这个结构储存的方块实体
	Block nbt_parser_interface.Block
}

// Hash 给出该结构所包含方块实体唯一哈希校验和。
// 对于两个完全相同的 NBT 方块，它们应当具有一致
// 的哈希校验和
func (s StructureNBTBlock) Hash() uint64 {
	return nbt_hash.NBTBlockHash(s.Block)
}
