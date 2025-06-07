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
			return index, nearBlockMapping[index], value
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
			return index, nearBlockMapping[index], value
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
			return index, nearBlockMapping[index], value
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
			return index, nearBlockMapping[index], value
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
			return index, nearBlockMapping[index], value
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
			return index, nearBlockMapping[index], value
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
// 如果 inclueEmptyContainer 为真，则会优先
// 考虑已经是全空的容器，否则优先考虑空气方块。
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
	if inclueEmptyContainer {
		index, offset, block = c.FindEmptyContainer(includeCenter)
		if block != nil {
			return
		}
	}

	index, offset, block = c.FindAir(includeCenter)
	if block != nil {
		return
	}

	index, offset, block = c.FindNonContainerAndNonAnvil(includeCenter)
	if block != nil {
		return
	}

	index, offset, block = c.FindNonAnvil(includeCenter)
	return
}

// FindSpaceToPlaceHelper 尝试从操作台找到一个位置以便于使用者
// 在该处放置帮助方块。帮助方块应当是铁砧或织布机。
//
// 对于容器，应当使用 FindSpaceToPlaceNewContainer 进行查找。
//
// includeCenter 指示要查找的方块是否也包括操作台中心处的方块。
//
// 如果 inclueEmptyContainer 为真，则会优先考虑已经是全空的容器，
// 否则优先考虑非铁砧方块。
//
// 返回的 index 可用于 BlockByIndex，
// 而返回的 offset 可用于 BlockByOffset。
//
// 如果返回的 block 不为空，则说明目标方块被找到，否则没有找到。
// 找到的方块可以通过修改其指向的值从而将它变成其他方块
func (c Console) FindSpaceToPlaceHelper(includeCenter bool, inclueEmptyContainer bool) (
	index int,
	offset protocol.BlockPos,
	block *block_helper.BlockHelper,
) {
	if inclueEmptyContainer {
		index, offset, block = c.FindEmptyContainer(includeCenter)
		if block != nil {
			return
		}
	}

	index, offset, block = c.FindAir(includeCenter)
	if block != nil {
		return
	}

	index, offset, block = c.FindNonAnvil(includeCenter)
	return
}

// FindOrGenerateNewAnvil 寻找操作台的 4 个帮助方块中
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

	block = c.NearBlockByIndex(index, protocol.BlockPos{0, -1, 0})
	if _, ok := (*block).(block_helper.Air); ok {
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
