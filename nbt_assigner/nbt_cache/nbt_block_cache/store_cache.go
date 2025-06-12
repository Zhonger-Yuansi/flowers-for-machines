package nbt_block_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
)

// StoreCache 将操作台中心处的 NBT 方块保存到当前的缓存命中系统。
// block 是操作台中心处的 NBT 方块数据；
// offset 是相对于这个 NBT 方块的偏移，例如床尾相对于床头的偏移
func (n *NBTBlockCache) StoreCache(block nbt_parser_interface.Block, offset protocol.BlockPos) error {
	var err error

	structure := StructureNBTBlock{
		HashNumber: nbt_hash.CompletelyHashNumber{
			HashNumber:    nbt_hash.NBTBlockFullHash(block),
			SetHashNumber: nbt_hash.ContainerSetHash(block),
		},
		Offset: offset,
		Block:  block,
	}

	_, ok := n.cachedNBTBlock[structure.HashNumber.HashNumber]
	if ok {
		return nil
	}

	structure.UniqueID, err = n.console.API().StructureBackup().BackupOffset(
		n.console.Center(),
		structure.Offset,
	)
	if err != nil {
		return fmt.Errorf("StoreCache: %v", err)
	}

	n.cachedNBTBlock[structure.HashNumber.HashNumber] = structure
	return nil
}

// CleanCache 清除该缓存命中系统中已有的全部缓存
func (n *NBTBlockCache) CleanCache() {
	api := n.console.API().StructureBackup()

	for _, value := range n.cachedNBTBlock {
		_ = api.DeleteStructure(value.UniqueID)
	}

	n.cachedNBTBlock = make(map[uint64]StructureNBTBlock)
}
