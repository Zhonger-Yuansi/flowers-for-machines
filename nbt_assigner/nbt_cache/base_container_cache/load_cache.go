package base_container_cache

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// LoadCache 加载名为 name 且方块状态为 states 的基容器。
// customName 指示基容器的自定义名称，它通常情况下为空。
// 如果没有找到，则尝试从已保存的结构中加载。
// 应当说明的是，基容器会在操作台的中心处被加载
func (b *BaseContainerCache) LoadCache(name string, states map[string]any, customName string) (hit bool, err error) {
	// Compute hash number
	container := BaseContainer{
		BlockName:         name,
		BlockStatesString: utils.MarshalBlockStates(states),
	}
	hashNumber := container.Hash()

	// Try to find target container from console
	block := b.console.BlockByIndex(nbt_console.ConsoleIndexCenterBlock)
	c, ok := (*block).(block_helper.ContainerBlockHelper)
	if ok && c.IsEmpty {
		currentContainer := BaseContainer{
			BlockName:         c.BlockName(),
			BlockStatesString: utils.MarshalBlockStates(c.BlockStates()),
			CustomeName:       customName,
		}
		if currentContainer.Hash() == hashNumber {
			return true, nil
		}
	}

	// Try to load from internal structure record mapping
	structure, ok := b.cachedBaseContainer[hashNumber]
	if !ok {
		return false, nil
	}

	// Load structure
	err = b.console.API().StructureBackup().RevertStructure(
		structure.UniqueID,
		b.console.Center(),
	)
	if err != nil {
		return false, nil
	}

	// Update underlying container data
	newContainer := block_helper.ContainerBlockHelper{
		OpenInfo: structure.Container,
		IsEmpty:  true,
	}
	b.console.UseHelperBlock(b.uniqueID, nbt_console.ConsoleIndexCenterBlock, newContainer)

	// Return
	return true, nil
}
