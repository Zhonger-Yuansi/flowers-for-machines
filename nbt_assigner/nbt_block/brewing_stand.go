package nbt_block

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
	nbt_parser_item "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/item"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// 酿造台
type BrewingStand struct {
	console *nbt_console.Console
	data    nbt_parser_block.BrewingStand
}

func (BrewingStand) Offset() protocol.BlockPos {
	return protocol.BlockPos{0, 0, 0}
}

func (b *BrewingStand) Make() error {
	api := b.console.API()
	usedSyncReplaceitemCommand := false
	existItemNeedRename := false

	brewingStandStates := map[string]any{
		"brewing_stand_slot_a_bit": byte(0),
		"brewing_stand_slot_b_bit": byte(0),
		"brewing_stand_slot_c_bit": byte(0),
	}
	updateBlockStates := func() {
		b.console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.ContainerBlockHelper{
			OpenInfo: block_helper.ContainerBlockOpenInfo{
				Name:                  b.data.BlockName(),
				States:                brewingStandStates,
				ConsiderOpenDirection: false,
			},
		})
	}

	// 清除酿造台处的方块
	err := api.SetBlock().SetBlock(b.console.Center(), "minecraft:air", "[]")
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}
	b.console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.Air{})

	// 生成酿造台方块
	if len(b.data.NBT.CustomName) == 0 {
		// 这个酿造台没有自定义物品名称，
		// 可以直接采用 Setblock 放置
		err = api.SetBlock().SetBlock(b.console.Center(), b.data.BlockName(), utils.MarshalBlockStates(brewingStandStates))
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	} else {
		// 这个酿造台具有自定义物品名称，
		// 需要进一步特殊处理
		err = b.console.API().Replaceitem().ReplaceitemInInventory(
			"@s",
			game_interface.ReplacePathHotbarOnly,
			game_interface.ReplaceitemInfo{
				Name:     "minecraft:brewing_stand",
				Count:    1,
				MetaData: 0,
				Slot:     b.console.HotbarSlotID(),
			},
			"",
			false,
		)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		b.console.UseInventorySlot(nbt_console.RequesterUser, b.console.HotbarSlotID(), true)
		// 打开铁砧
		index, err := b.console.FindOrGenerateNewAnvil()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		success, err := b.console.OpenContainerByIndex(index)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			return fmt.Errorf("Make: Failed to open the anvil in setblock stage")
		}
		// 物品改名
		success, _, _, err = api.ItemStackOperation().OpenTransaction().
			RenameInventoryItem(b.console.HotbarSlotID(), b.data.NBT.CustomName).
			Commit()
		if err != nil {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: %v", err)
		}
		if !success {
			_ = api.ContainerOpenAndClose().CloseContainer()
			return fmt.Errorf("Make: The server rejected the renaming operation (setblock stage)")
		}
		// 关闭铁砧
		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		// 前往操作台中心处
		err = b.console.CanReachOrMove(b.console.Center())
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		// 点击地板以放置酿造台
		_, _, err = api.BotClick().PlaceBlockHighLevel(b.console.Center(), b.console.HotbarSlotID(), 1)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		*b.console.NearBlockByIndex(nbt_console.ConsoleIndexCenterBlock, protocol.BlockPos{0, -1, 0}) = block_helper.NearBlock{
			Name: game_interface.BasePlaceBlock,
		}
	}
	updateBlockStates()

	// 处理可以直接 Replaceitem 处理的物品
	for _, item := range b.data.NBT.Items {
		var fuelAddCount uint8 = 0
		underlaying := item.Item.UnderlyingItem()
		defaultItem := underlaying.(*nbt_parser_item.DefaultItem)

		if item.Item.NeedEnchOrRename() {
			existItemNeedRename = true
			continue
		}
		if item.Slot == 4 {
			fuelAddCount = 1
		}

		usedSyncReplaceitemCommand = true
		switch item.Slot {
		case 1:
			brewingStandStates["brewing_stand_slot_a_bit"] = byte(1)
		case 2:
			brewingStandStates["brewing_stand_slot_b_bit"] = byte(1)
		case 3:
			brewingStandStates["brewing_stand_slot_c_bit"] = byte(1)
		}

		err = b.console.API().Replaceitem().ReplaceitemInContainerAsync(
			b.console.Center(),
			game_interface.ReplaceitemInfo{
				Name:     item.Item.ItemName(),
				Count:    item.Item.ItemCount() + fuelAddCount,
				MetaData: item.Item.ItemMetadata(),
				Slot:     resources_control.SlotID(item.Slot),
			},
			utils.MarshalItemComponent(defaultItem.Enhance.ItemComponent),
		)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		updateBlockStates()
	}

	// 如果使用了 Replaceitem 命令，
	// 则需要等待更改
	if usedSyncReplaceitemCommand {
		err = api.Commands().AwaitChangesGeneral()
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// 如果没有物品存在自定义物品名称，
	// 则可以直接返回值
	if !existItemNeedRename {
		return nil
	}

	// 先将需要特殊处理的物品放入快捷栏
	for _, item := range b.data.NBT.Items {
		var fuelAddCount uint8 = 0
		underlaying := item.Item.UnderlyingItem()
		defaultItem := underlaying.(*nbt_parser_item.DefaultItem)

		if !item.Item.NeedEnchOrRename() {
			continue
		}
		if item.Slot == 4 {
			fuelAddCount = 1
		}

		err = b.console.API().Replaceitem().ReplaceitemInInventory(
			"@s",
			game_interface.ReplacePathHotbarOnly,
			game_interface.ReplaceitemInfo{
				Name:     item.Item.ItemName(),
				Count:    item.Item.ItemCount() + fuelAddCount,
				MetaData: item.Item.ItemMetadata(),
				Slot:     resources_control.SlotID(item.Slot),
			},
			utils.MarshalItemComponent(defaultItem.Enhance.ItemComponent),
			false,
		)
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
		b.console.UseInventorySlot(nbt_console.RequesterUser, resources_control.SlotID(item.Slot), true)
	}

	// 打开铁砧
	index, err := b.console.FindOrGenerateNewAnvil()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}
	success, err := b.console.OpenContainerByIndex(index)
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}
	if !success {
		return fmt.Errorf("Make: Failed to open the anvil in brewing stand item rename stage")
	}

	// 物品重命名
	transaction := api.ItemStackOperation().OpenTransaction()
	for _, item := range b.data.NBT.Items {
		underlaying := item.Item.UnderlyingItem()
		defaultItem := underlaying.(*nbt_parser_item.DefaultItem)
		if !item.Item.NeedEnchOrRename() {
			continue
		}
		_ = transaction.RenameInventoryItem(resources_control.SlotID(item.Slot), defaultItem.Enhance.DisplayName)
	}

	// 提交更改
	success, _, _, err = transaction.Commit()
	if err != nil {
		_ = api.ContainerOpenAndClose().CloseContainer()
		return fmt.Errorf("Make: %v", err)
	}
	if !success {
		_ = api.ContainerOpenAndClose().CloseContainer()
		return fmt.Errorf("Make: The server rejected the renaming operation (brewing stand item rename stage)")
	}

	// 关闭铁砧
	err = api.ContainerOpenAndClose().CloseContainer()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	// 打开酿造台
	success, err = b.console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}
	if !success {
		return fmt.Errorf("Make: Failed to open the brewing stand")
	}
	defer api.ContainerOpenAndClose().CloseContainer()

	// 移动已改名物品到酿造台
	for _, item := range b.data.NBT.Items {
		var fuelAddCount uint8 = 0
		if !item.Item.NeedEnchOrRename() {
			continue
		}

		switch item.Slot {
		case 1:
			brewingStandStates["brewing_stand_slot_a_bit"] = byte(1)
		case 2:
			brewingStandStates["brewing_stand_slot_b_bit"] = byte(1)
		case 3:
			brewingStandStates["brewing_stand_slot_c_bit"] = byte(1)
		case 4:
			fuelAddCount = 1
		}

		_ = transaction.MoveToContainer(
			resources_control.SlotID(item.Slot),
			resources_control.SlotID(item.Slot),
			item.Item.ItemCount()+fuelAddCount,
		)
	}

	// 提交更改
	success, _, _, err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}
	if !success {
		return fmt.Errorf("Make: The server rejected the stack request action")
	}

	// 更新方块状态
	brewingStandStates = b.data.BlockStates()
	updateBlockStates()

	// 返回值
	return nil
}
