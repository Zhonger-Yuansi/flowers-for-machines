package nbt_console

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
)

// FindAir 从操作台的帮助方块中寻找一个空气方块。
// includeCenter 指示要查找的方块是否也包括操作台
// 中心处的方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明找到，
// 否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindAir(includeCenter bool) (index int, offset protocol.BlockPos, block *block_helper.BlockHelper) {
	for index, value := range c.helperBlocks {
		if !includeCenter && index == 0 {
			continue
		}
		if _, ok := (*value).(block_helper.Air); ok {
			return index, helperBlockMapping[index], value
		}
	}
	return 0, protocol.BlockPos{}, nil
}

// FindAir 从操作台的帮助方块中寻找一个铁砧方块。
// includeCenter 指示要查找的方块是否也包括操作
// 台中心处的方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明找到，
// 否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindAnvil(includeCenter bool) (index int, offset protocol.BlockPos, block *block_helper.BlockHelper) {
	for index, value := range c.helperBlocks {
		if !includeCenter && index == 0 {
			continue
		}
		if _, ok := (*value).(block_helper.AnvilBlockHelper); ok {
			return index, helperBlockMapping[index], value
		}
	}
	return 0, protocol.BlockPos{}, nil
}

// FindLoom 从操作台的帮助方块中寻找一个织布机方块。
// includeCenter 指示要查找的方块是否也包括操作台
// 中心处的方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明找到，
// 否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindLoom(includeCenter bool) (index int, offset protocol.BlockPos, block *block_helper.BlockHelper) {
	for index, value := range c.helperBlocks {
		if !includeCenter && index == 0 {
			continue
		}
		if _, ok := (*value).(block_helper.LoomBlockHelper); ok {
			return index, helperBlockMapping[index], value
		}
	}
	return 0, protocol.BlockPos{}, nil
}

// FindNonAnvil 从操作台的帮助方块中寻找一
// 个非铁砧方块。这意味目标方块将可以是空气、
// 容器，或织布机。
//
// includeCenter 指示要查找的方块是否也包括操
// 作台中心处的方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明找到，
// 否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindNonAnvil(includeCenter bool) (index int, offset protocol.BlockPos, block *block_helper.BlockHelper) {
	for index, value := range c.helperBlocks {
		if !includeCenter && index == 0 {
			continue
		}
		if _, ok := (*value).(block_helper.AnvilBlockHelper); !ok {
			return index, helperBlockMapping[index], value
		}
	}
	return 0, protocol.BlockPos{}, nil
}

// FindNonContainerAndNonAnvil 从操作台的帮
// 助方块中寻找一个不是容器且也不是铁砧的方块。
// 这意味目标方块将可以是空气或织布机。
//
// includeCenter 指示要查找的方块是否也包括操
// 作台中心处的方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明找到，
// 否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindNonContainerAndNonAnvil(includeCenter bool) (index int, offset protocol.BlockPos, block *block_helper.BlockHelper) {
	for index, value := range c.helperBlocks {
		if !includeCenter && index == 0 {
			continue
		}
		switch (*value).(type) {
		case block_helper.ContainerBlockHelper:
		case block_helper.AnvilBlockHelper:
		default:
			return index, helperBlockMapping[index], value
		}
	}
	return 0, protocol.BlockPos{}, nil
}

// FindEmptyContainer 从操作台的帮助方块
// 中寻找一个全空的容器方块。
//
// includeCenter 指示要查找的方块是否也包
// 括操作台中心处的方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明找到，
// 否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindEmptyContainer(includeCenter bool) (index int, offset protocol.BlockPos, block *block_helper.BlockHelper) {
	for index, value := range c.helperBlocks {
		if !includeCenter && index == 0 {
			continue
		}
		if container, ok := (*value).(block_helper.ContainerBlockHelper); ok && container.IsEmpty {
			return index, helperBlockMapping[index], value
		}
	}
	return 0, protocol.BlockPos{}, nil
}

// FindSpaceToPlaceNewContainer 尝试从操作台
// 找到一个位置以便于使用者在该处放置新容器。
//
// includeCenter 指示要查找的方块是否也包括操
// 作台中心处的方块。
//
// 如果 inclueEmptyContainer 为真，则会在优先
// 考虑空气后优先考虑空容器。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明目标方块被
// 找到，否则没有找到。找到的方块可以通过修改
// 其指向的值从而将它变成其他方块
func (c Console) FindSpaceToPlaceNewContainer(includeCenter bool, inclueEmptyContainer bool) (
	index int,
	offset protocol.BlockPos,
	block *block_helper.BlockHelper,
) {
	index, offset, block = c.FindAir(includeCenter)
	if block != nil {
		return
	}

	if inclueEmptyContainer {
		index, offset, block = c.FindEmptyContainer(includeCenter)
		if block != nil {
			return
		}
	}

	index, offset, block = c.FindNonContainerAndNonAnvil(includeCenter)
	return
}

// FindSpaceToPlaceHelper 尝试从操作台找到一个位置以便于使用者
// 在该处放置帮助方块。帮助方块应当是铁砧或织布机。
//
// 对于容器，应当使用 FindSpaceToPlaceNewContainer 进行查找。
// includeCenter 指示要查找的方块是否也包括操作台中心处的方块。
//
// 如果 inclueEmptyContainer 为真，
// 则会在优先考虑空气后优先考虑空容器。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 可以保证在正确使用操作台的情况下一定可以找到需要的方块。
// 另外，找到的方块可以通过修改其指向的值从而将它变成其他方块
func (c Console) FindSpaceToPlaceHelper(includeCenter bool, inclueEmptyContainer bool) (
	index int,
	offset protocol.BlockPos,
	block *block_helper.BlockHelper,
) {
	index, offset, block = c.FindAir(includeCenter)
	if block != nil {
		return
	}

	if inclueEmptyContainer {
		index, offset, block = c.FindEmptyContainer(includeCenter)
		if block != nil {
			return
		}
	}

	index, offset, block = c.FindNonAnvil(includeCenter)
	if block == nil {
		panic("FindSpaceToPlaceHelper: Should nerver happend")
	}

	return
}

// FindOrGenerateNewAnvil 寻找操作台的 8 个帮助方块中
// 是否有一个是铁砧。如果没有，则生成一个铁砧及其承重方块。
// index 指示找到或生成的铁砧在操作台上的索引
func (c *Console) FindOrGenerateNewAnvil() (index int, err error) {
	var block *block_helper.BlockHelper
	var needFloorBlock bool

	index, _, block = c.FindAnvil(false)
	if block != nil {
		return
	}

	index, _, block = c.FindSpaceToPlaceHelper(false, false)
	if block == nil {
		panic("FindOrGenerateNewAnvil: Should nerver happened")
	}

	nearBlock := c.NearBlockByIndex(index, protocol.BlockPos{0, -1, 0})
	if _, ok := (*nearBlock).(block_helper.Air); ok {
		needFloorBlock = true
	}

	states, err := c.api.SetBlock().SetAnvil(c.BlockPosByIndex(index), needFloorBlock)
	if err != nil {
		return 0, fmt.Errorf("FindOrGenerateNewAnvil: %v", err)
	}

	anvil := block_helper.AnvilBlockHelper{States: states}
	c.UseHelperBlock(RequesterSystemCall, index, anvil)
	if needFloorBlock {
		var floorBlock block_helper.BlockHelper = block_helper.NearBlock{
			Name: game_interface.BaseAnvil,
		}
		*c.NearBlockByIndex(index, protocol.BlockPos{0, -1, 0}) = floorBlock
	}

	return index, nil
}

// FindOrGenerateNewLoom 寻找操作台的 8 个帮助方块中
// 是否有一个是织布机。如果没有，则生成一个新的织布机。
// index 指示找到或生成的织布机在操作台上的索引
func (c *Console) FindOrGenerateNewLoom() (index int, err error) {
	var block *block_helper.BlockHelper

	index, _, block = c.FindLoom(false)
	if block != nil {
		return
	}

	index, _, block = c.FindSpaceToPlaceHelper(false, false)
	if block == nil {
		panic("FindOrGenerateNewLoom: Should nerver happened")
	}

	loom := block_helper.LoomBlockHelper{}
	err = c.api.SetBlock().SetBlock(
		c.BlockPosByIndex(index),
		loom.BlockName(),
		loom.BlockStatesString(),
	)
	if err != nil {
		return 0, fmt.Errorf("FindOrGenerateNewLoom: %v", err)
	}
	c.UseHelperBlock(RequesterSystemCall, index, loom)

	return index, nil
}
