package nbt_assigner

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
)

const (
	// BaseBackground 是操作台地板的构成方块
	BaseBackground = "verdant_froglight"
	// DefaultTimeoutInitConsole 是抵达操作台目标区域的最长等待期限
	DefaultTimeoutInitConsole = time.Second * 30
)

var (
	// offsetMapping ..
	offsetMapping = []protocol.BlockPos{
		[3]int32{-1, 0, 0},
		[3]int32{1, 0, 0},
		[3]int32{0, -1, 0},
		[3]int32{0, 1, 0},
		[3]int32{0, 0, 1},
		[3]int32{0, 0, -1},
	}
	// offsetMappingInv ..
	offsetMappingInv = map[protocol.BlockPos]int{
		[3]int32{-1, 0, 0}: 0,
		[3]int32{1, 0, 0}:  1,
		[3]int32{0, -1, 0}: 2,
		[3]int32{0, 1, 0}:  3,
		[3]int32{0, 0, 1}:  4,
		[3]int32{0, 0, -1}: 5,
	}
	// nearBlockMapping ..
	nearBlockMapping = []protocol.BlockPos{
		[3]int32{0, 0, 0},
		[3]int32{-3, 0, 0},
		[3]int32{3, 0, 0},
		[3]int32{0, 0, 3},
		[3]int32{0, 0, -3},
	}
	// nearBlockMappingInv ..
	nearBlockMappingInv = map[protocol.BlockPos]int{
		[3]int32{0, 0, 0}:  0,
		[3]int32{-3, 0, 0}: 1,
		[3]int32{3, 0, 0}:  2,
		[3]int32{0, 0, 3}:  3,
		[3]int32{0, 0, -3}: 4,
	}
)

// Console 是机器人导入 NBT 方块所使用的操作台。
// 它目前被定义为一个 11*5*11 的全空气区域
type Console struct {
	// api 是与租赁服进行交互的若干接口
	api *game_interface.GameInterface
	// center 是操作台的中心位置
	center protocol.BlockPos
	// position 是机器人目前所在的方块位置
	position protocol.BlockPos
	// helperBlocks 是操作台中心及其
	// 不远处等距离分布的 4 个帮助方块。
	// 通过记录这 5 个方块的实际情况，
	// 有助于减少部分操作的实际耗时
	helperBlocks [5]*block_helper.BlockHelper
	// nearBlocks 是操作台中心方块及另
	// 外 4 个帮助方块相邻的方块。
	//
	// 如果认为操作台中心方块和另外 4 个
	// 帮助方块是 master 方块，那么对于
	// 第二层数组，可以通过 offsetMapping
	// 确定它们各自相邻其 master 方块的位置变化。
	//
	// 另外，offsetMappingInv 是 offsetMapping
	// 的逆映射
	nearBlocks [5][6]*block_helper.BlockHelper
}

// NewConsole 根据交互接口 api 和操作台中心 center 创建并返回一个新的操作台实例
func NewConsole(api *game_interface.GameInterface, center protocol.BlockPos) *Console {
	c := &Console{
		api:    api,
		center: center,
	}

	for index := range 5 {
		var airBlock block_helper.BlockHelper = block_helper.AirBlock{}
		c.helperBlocks[index] = &airBlock
	}
	for index := range 5 {
		for idx := range 6 {
			var airBlock block_helper.BlockHelper = block_helper.AirBlock{}
			c.nearBlocks[index][idx] = &airBlock
		}
	}

	return c
}

// InitConsoleArea 将机器人传送至操作台的中心方块处，
// 并试图初始化操作台的地板方块。
//
// InitConsoleArea 应当至多调用一次，并且应在 NewConsole
// 尽可能快的调用。
//
// InitConsoleArea 的调用者有责任确保操作台位于主世界，
// 并且操作台中心方块处的 11*5*11 的区域全为空气且没有
// 任何实体
func (c *Console) InitConsoleArea() error {
	api := c.api.Commands()

	err := api.SendSettingsCommand(
		fmt.Sprintf("execute in overworld run tp %d %d %d", c.center[0], c.center[1], c.center[2]),
		true,
	)
	if err != nil {
		return fmt.Errorf("InitConsoleArea: %v", err)
	}

	timer := time.NewTimer(DefaultTimeoutInitConsole)
	defer timer.Stop()

	for {
		resp, err := api.SendWSCommandWithResp(
			fmt.Sprintf(
				"execute as @s at @s positioned %d ~ ~ positioned ~ 0 ~ positioned ~ ~ %d run testforblock ~ 320 ~ air",
				c.center[0], c.center[2],
			),
		)
		if err != nil {
			return fmt.Errorf("InitConsoleArea: %v", err)
		}

		if resp.SuccessCount > 0 {
			c.position = c.center
			break
		}

		select {
		case <-timer.C:
			return fmt.Errorf("InitConsoleArea: Can not teleport to the target area (timeout)")
		default:
		}
	}

	_, err = api.SendWSCommandWithResp(
		fmt.Sprintf(
			"execute as @s at @s positioned %d ~ ~ positioned ~ %d ~ positioned ~ ~ %d run fill ~-5 ~-2 ~-5 ~5 ~2 ~5 air",
			c.center[0], c.center[1], c.center[2],
		),
	)
	if err != nil {
		return fmt.Errorf("InitConsoleArea: %v", err)
	}

	err = api.SendSettingsCommand(
		fmt.Sprintf("execute in overworld run tp %d %d %d", c.center[0], c.center[1], c.center[2]),
		true,
	)
	if err != nil {
		return fmt.Errorf("InitConsoleArea: %v", err)
	}

	_, err = api.SendWSCommandWithResp(
		fmt.Sprintf(
			"execute as @s at @s positioned %d ~ ~ positioned ~ %d ~ positioned ~ ~ %d run fill ~-5 ~-1 ~-5 ~5 ~-1 ~5 %s",
			c.center[0], c.center[1], c.center[2], BaseBackground,
		),
	)
	if err != nil {
		return fmt.Errorf("InitConsoleArea: %v", err)
	}

	for index := range 5 {
		var baseBlock block_helper.BlockHelper = block_helper.NearBlockHelper{
			Name: BaseBackground,
		}
		c.nearBlocks[index][offsetMappingInv[[3]int32{0, -1, 0}]] = &baseBlock
	}

	return nil
}

// BlockByIndex 按 index 查找操作台上的方块。
// index 为 0 将查找操作台中心的方块，
// index 为 i (i>0) 将查找第 i-1 个帮助方块
func (c Console) BlockByIndex(index int) (result *block_helper.BlockHelper) {
	return c.helperBlocks[index]
}

// BlockByOffset 按坐标偏移量查找操作台上的方块。
// 如果给出的偏移量不能对应操作台上的方块，则返回操作台中心处的方块
func (c Console) BlockByOffset(offset protocol.BlockPos) (result *block_helper.BlockHelper) {
	return c.helperBlocks[nearBlockMappingInv[offset]]
}

// NearBlockByIndex 按 index 查找操作台上方块的相邻方块。
// index 为 0 将查找操作台中心处方块的相邻，
// index 为 i (i>0) 将查找第 i-1 个帮助方块的相邻。
//
// offset 指示在根据 index 找到目标方块后，相邻其的方块相对于
// 目标方块本身的坐标偏移。然后，我们偏移到该方块上并返回该方块
func (c Console) NearBlockByIndex(index int, offset protocol.BlockPos) (result *block_helper.BlockHelper) {
	return c.nearBlocks[index][offsetMappingInv[offset]]
}

// NearBlockByOffset 按 offset 查找操作台上方块的相邻。
// 如果给出的偏移量不能对应操作台上的方块，则返回操作台中心处方块的相邻。
//
// nearOffset 指示在根据 index 找到目标方块后，相邻其的方块相对于目标
// 方块本身的坐标偏移。然后，我们偏移到该方块上并返回该方块
func (c Console) NearBlockByOffset(offset protocol.BlockPos, nearOffset protocol.BlockPos) (result *block_helper.BlockHelper) {
	return c.nearBlocks[nearBlockMappingInv[offset]][offsetMappingInv[offset]]
}
