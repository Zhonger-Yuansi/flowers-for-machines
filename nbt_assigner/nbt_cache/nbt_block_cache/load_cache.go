package nbt_block_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
)

// LoadCache 尝试加载一个已缓存的 NBT 方块到操作台中心。
// 如果 hashNumber 所指示的缓存不存在，则不执行任何操作
func (n *NBTBlockCache) LoadCache(hashNumber uint64) error {
	structure, ok := n.cachedNBTBlock[hashNumber]
	if !ok {
		return nil
	}

	err := n.console.API().StructureBackup().RevertStructure(
		structure.UniqueID,
		n.console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
	)
	if err != nil {
		return fmt.Errorf("LoadCache: %v", err)
	}

	container, ok := structure.Block.(*nbt_parser_block.Container)
	if ok {
		n.console.UseHelperBlock(
			n.uniqueID,
			nbt_console.ConsoleIndexCenterBlock,
			block_helper.ContainerBlockHelper{
				OpenInfo: block_helper.ContainerBlockOpenInfo{
					Name:                  container.BlockName(),
					States:                container.BlockStates(),
					ConsiderOpenDirection: container.ConsiderOpenDirection(),
					ShulkerFacing:         container.NBT.ShulkerFacing,
				},
				IsEmpty: len(container.NBT.Items) == 0,
			},
		)
		return nil
	}

	n.console.UseHelperBlock(
		n.uniqueID,
		nbt_console.ConsoleIndexCenterBlock,
		block_helper.ComplexBlock{
			Name:   structure.Block.BlockName(),
			States: structure.Block.BlockStates(),
		},
	)
	return nil
}
