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
	// ..
	Block nbt_parser_interface.Block
}

// // Hash 给出这个基容器的唯一哈希校验和。
// // 校验和不包括该容器所在结构的唯一标识，
// // 这意味着来自两个不同结构的相同基容器
// // 具有完全相同的哈希校验和
// func (b BaseContainer) Hash() uint64 {
// 	buf := bytes.NewBuffer(nil)
// 	w := protocol.NewWriter(buf, 0)

// 	length := uint32(len(b.BlockName))
// 	w.Varuint32(&length)
// 	w.String(&b.BlockName)
// 	length = uint32(len(b.BlockStatesString))
// 	w.Varuint32(&length)
// 	w.String(&b.BlockStatesString)

// 	return xxhash.Sum64(buf.Bytes())
// }
