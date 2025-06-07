package base_container_cache

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// StoreCache 将操作台中心处的方块保存到当前缓存命中系统。
// 如果该方块不是一个容器，或其中装有物品，则返回错误
func (b *BaseContainerCache) StoreCache() error {
	block := b.console.BlockByIndex(nbt_console.ConsoleIndexCenterBlock)
	container, ok := (*block).(block_helper.ContainerBlockHelper)
	if !ok {
		return fmt.Errorf("StoreCache: The center of the console is not a container; *block = %#v", *block)
	}
	if !container.IsEmpty {
		return fmt.Errorf("StoreCache: Target container is not empty; *block = %#v", *block)
	}

	c := BaseContainer{
		BlockName:         container.OpenInfo.Name,
		BlockStatesString: utils.MarshalBlockStates(container.OpenInfo.States),
	}
	hashNumber := c.Hash()

	if _, ok := b.cachedBaseContainer[hashNumber]; ok {
		return nil
	}

	uniqueID, err := b.console.API().StructureBackup().BackupStructure(
		b.console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
	)
	if err != nil {
		return fmt.Errorf("StoreCache: %v", err)
	}

	b.cachedBaseContainer[hashNumber] = StructureBaseContainer{
		UniqueID:  uniqueID,
		Container: container.OpenInfo,
	}
	return nil
}
