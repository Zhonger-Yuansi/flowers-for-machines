package base_container_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// LoadCache 试图从操作台上查找一个名称为 name 且方块状态为 states 的基容器。
// 如果没有找到，则尝试从已经保存的结构中查找，然后将其加载到操作台上。
// 返回的 index 指示找到的方块在操作台上的索引
func (b *BaseContainerCache) LoadCache(name string, states map[string]any) (index int, hit bool, err error) {
	// Compute hash number
	container := BaseContainer{
		BlockName:         name,
		BlockStatesString: utils.MarshalBlockStates(states),
	}
	hashNumber := container.Hash()

	// Try to find target container from console
	for index := range 9 {
		if index == nbt_console.ConsoleIndexCenterBlock {
			continue
		}

		block := b.console.BlockByIndex(index)
		c, ok := (*block).(block_helper.ContainerBlockHelper)
		if !ok || !c.IsEmpty {
			continue
		}

		currentContainer := BaseContainer{
			BlockName:         c.BlockName(),
			BlockStatesString: utils.MarshalBlockStates(c.BlockStates()),
		}
		if currentContainer.Hash() == hashNumber {
			return index, true, nil
		}
	}

	// Try to load from internal structure record mapping
	structure, ok := b.cachedBaseContainer[hashNumber]
	if !ok {
		return 0, false, nil
	}

	// Find a place to place the container
	index, _, block := b.console.FindSpaceToPlaceNewContainer(false, false)
	if block == nil {
		index = nbt_console.ConsoleIndexFirstHelperBlock
	}

	// Load structure
	err = b.console.API().StructureBackup().RevertStructure(
		structure.UniqueID,
		b.console.BlockPosByIndex(index),
	)
	if err != nil {
		return 0, false, fmt.Errorf("LoadCache: %v", err)
	}

	// Update underlying container data
	newContainer := block_helper.ContainerBlockHelper{
		OpenInfo: structure.Container,
		IsEmpty:  true,
	}
	b.console.UseHelperBlock(b.uniqueID, index, newContainer)

	// Return
	return index, true, nil
}
