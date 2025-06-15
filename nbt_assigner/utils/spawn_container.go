package nbt_assigner_utils

import (
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
)

// SpawnContainer 在操作台中心生成空的 container。
// 确保这个容器的自定义物品名称也被考虑在内。
//
// SpawnContainer 因为可能需要通过点击放置容器，
// 因此快捷栏的槽位可能会被重用。调用者有责任确保
// 快捷栏的物品不会因此而被意外使用
func SpawnContainer(
	console *nbt_console.Console,
	cache *nbt_cache.NBTCacheSystem,
	container nbt_parser_block.Container,
) error {
	// 准备
	api := console.API()
	useCommandToPlaceBlock := true
	successFunc := func() {
		console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.ContainerBlockHelper{
			OpenInfo: block_helper.ContainerBlockOpenInfo{
				Name:                  container.BlockName(),
				States:                container.BlockStates(),
				ConsiderOpenDirection: container.ConsiderOpenDirection(),
				ShulkerFacing:         container.NBT.ShulkerFacing,
			},
		})
	}

	// 尝试基容器缓存
	hit, err := cache.BaseContainerCache().LoadCache(
		container.BlockName(),
		container.BlockStates(),
		container.CustomName,
		container.NBT.ShulkerFacing,
	)
	if err != nil {
		return fmt.Errorf("SpawnContainer: %v", err)
	}
	if hit {
		return nil
	}

	// 先将目标位置替换为空气
	err = api.SetBlock().SetBlock(console.Center(), "minecraft:air", "[]")
	if err != nil {
		return fmt.Errorf("SpawnContainer: %v", err)
	}
	console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.Air{})

	// 检查是否需要复杂的工序来放置容器
	if len(container.CustomName) > 0 {
		useCommandToPlaceBlock = false
	}
	if strings.Contains(container.BlockName(), "shulker") {
		if container.NBT.ShulkerFacing != 1 {
			useCommandToPlaceBlock = false
		}
	}

	// 如果需要复杂的工序
	if !useCommandToPlaceBlock {
		// 先把目标物品获取到物品栏
		err := api.Replaceitem().ReplaceitemInInventory(
			"@s",
			game_interface.ReplacePathHotbarOnly,
			game_interface.ReplaceitemInfo{
				Name:     container.BlockName(),
				Count:    1,
				MetaData: 0,
				Slot:     console.HotbarSlotID(),
			},
			"",
			true,
		)
		if err != nil {
			return fmt.Errorf("SpawnContainer: %v", err)
		}
		console.UseInventorySlot(nbt_console.RequesterUser, console.HotbarSlotID(), true)

		// 这个容器具有自定义的物品名称，需要进一步特殊处理
		if len(container.CustomName) > 0 {
			index, err := console.FindOrGenerateNewAnvil()
			if err != nil {
				return fmt.Errorf("SpawnContainer: %v", err)
			}

			success, err := console.OpenContainerByIndex(index)
			if err != nil {
				return fmt.Errorf("SpawnContainer: %v", err)
			}
			if !success {
				return fmt.Errorf("SpawnContainer: Failed to open the anvil to rename container")
			}

			success, _, _, err = api.ItemStackOperation().OpenTransaction().
				RenameInventoryItem(console.HotbarSlotID(), container.CustomName).
				Commit()
			if err != nil {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("SpawnContainer: %v", err)
			}
			if !success {
				_ = api.ContainerOpenAndClose().CloseContainer()
				return fmt.Errorf("SpawnContainer: The server rejected the container rename operation")
			}

			err = api.ContainerOpenAndClose().CloseContainer()
			if err != nil {
				return fmt.Errorf("SpawnContainer: %v", err)
			}
		}

		// 确定放置目标容器时所使用的朝向
		var facing uint8 = 1
		if strings.Contains(container.BlockName(), "shulker") {
			facing = container.NBT.ShulkerFacing
		}

		// 移动机器人到操作台中心
		err = console.CanReachOrMove(console.Center())
		if err != nil {
			return fmt.Errorf("makeNormal: %v", err)
		}

		// 放置目标容器
		_, offsetPos, err := api.BotClick().PlaceBlockHighLevel(console.Center(), console.HotbarSlotID(), facing)
		if err != nil {
			return fmt.Errorf("makeNormal: %v", err)
		}
		successFunc()
		*console.NearBlockByIndex(nbt_console.ConsoleIndexCenterBlock, offsetPos) = block_helper.NearBlock{
			Name: game_interface.BasePlaceBlock,
		}

		// 将该容器保存到基容器缓存命中系统
		err = cache.BaseContainerCache().StoreCache(container.CustomName, container.NBT.ShulkerFacing)
		if err != nil {
			return fmt.Errorf("makeNormal: %v", err)
		}

		// 返回值
		return nil
	}

	// 目标容器可以直接通过简单的 setblock 放置
	err = api.SetBlock().SetBlock(console.Center(), container.BlockName(), container.BlockStatesString())
	if err != nil {
		return fmt.Errorf("SpawnContainer: %v", err)
	}
	successFunc()

	return nil
}
