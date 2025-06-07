package main

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache/item_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/pterm/pterm"
)

func SystemTestingItemCache() {
	tA := time.Now()

	// Test round 1
	{
		chestStatesString := `["minecraft:cardinal_direction"="east"]`

		api.Commands().SendSettingsCommand("gamemode 1", true)
		api.SetBlock().SetBlock(
			console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
			"chest",
			chestStatesString,
		)

		api.Replaceitem().ReplaceitemInContainerAsync(
			console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
			game_interface.ReplaceitemInfo{
				Name:     "apple",
				Count:    3,
				MetaData: 0,
				Slot:     2,
			},
			"",
		)
		api.Replaceitem().ReplaceitemInContainerAsync(
			console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
			game_interface.ReplaceitemInfo{
				Name:     "diamond_sword",
				Count:    1,
				MetaData: 7,
				Slot:     1,
			},
			"",
		)
		api.Commands().AwaitChangesGeneral()

		container := block_helper.ContainerBlockHelper{
			OpenInfo: block_helper.ContainerBlockOpenInfo{
				Name:                  "chest",
				States:                utils.ParseBlockStatesString(chestStatesString),
				ConsiderOpenDirection: true,
			},
			IsEmpty: false,
		}
		*console.BlockByIndex(nbt_console.ConsoleIndexCenterBlock) = container

		err := itemCache.StoreCache(
			[]item_cache.ItemCacheInfo{
				{SlotID: 2, Count: 3, Hash: item_cache.ItemHashNumber{HashNumber: 1, SetHashNumber: item_cache.SetHashNumberNotExist}},
				{SlotID: 1, Count: 1, Hash: item_cache.ItemHashNumber{HashNumber: 2, SetHashNumber: 1}},
			},
			container.OpenInfo,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 1 failed due to %v (stage 1)", err))
		}

		err = itemCache.StoreCache(
			[]item_cache.ItemCacheInfo{
				{SlotID: 1, Count: 1, Hash: item_cache.ItemHashNumber{HashNumber: 2, SetHashNumber: 1}},
			},
			container.OpenInfo,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 1 failed due to %v (stage 2)", err))
		}
	}

	// Test round 2
	{
		slotIDA, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    1,
				SetHashNumber: item_cache.SetHashNumberNotExist,
			},
			nil,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 2 failed due to %v (stage 1)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}

		slotIDB, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    1,
				SetHashNumber: item_cache.SetHashNumberNotExist,
			},
			nil,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 2 failed due to %v (stage 2)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if slotIDA != slotIDB {
			panic("SystemTestingItemCache: Test round 2 failed")
		}

		slotIDC, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    2,
				SetHashNumber: 1,
			},
			[]resources_control.SlotID{slotIDA},
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 2 failed due to %v (stage 3)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if slotIDC == slotIDA {
			panic("SystemTestingItemCache: Test round 2 failed")
		}

		api.Commands().SendWSCommandWithResp("clear")
		console.CleanInventory()

		slotIDD, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    2018,
				SetHashNumber: 1,
			},
			nil,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 2 failed due to %v (stage 4)", err))
		}
		if !isSetHashHit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if slotIDD != slotIDA {
			panic("SystemTestingItemCache: Test round 2 failed")
		}

		api.Replaceitem().ReplaceitemInInventory(
			"@s",
			game_interface.ReplacePathInventory,
			game_interface.ReplaceitemInfo{
				Name:     "air",
				Count:    1,
				MetaData: 0,
				Slot:     slotIDD,
			},
			"",
			true,
		)
		itemCache.ConsumeCache(slotIDD)

		slotIDE, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    1,
				SetHashNumber: item_cache.SetHashNumberNotExist,
			},
			nil,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 2 failed due to %v (stage 5)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
		if slotIDE != slotIDA {
			panic("SystemTestingItemCache: Test round 2 failed")
		}
	}

	// Test round 3
	{
		api.Commands().SendWSCommandWithResp("clear")
		console.CleanInventory()

		api.Replaceitem().ReplaceitemInInventory(
			"@s",
			game_interface.ReplacePathHotbarOnly,
			game_interface.ReplaceitemInfo{
				Name:     "light_blue_shulker_box",
				Count:    1,
				MetaData: 0,
				Slot:     console.HotbarSlotID(),
			},
			"",
			true,
		)

		_, offset, _ := api.BotClick().PlaceBlockHighLevel(
			console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
			console.HotbarSlotID(),
			0,
		)
		console.UpdatePosition(console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock))
		var nearBlock block_helper.BlockHelper = block_helper.NearBlock{
			Name: game_interface.BasePlaceBlock,
		}
		*console.NearBlockByIndex(nbt_console.ConsoleIndexCenterBlock, offset) = nearBlock
		api.Commands().AwaitChangesGeneral()

		api.Replaceitem().ReplaceitemInContainerAsync(
			console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
			game_interface.ReplaceitemInfo{
				Name:     "apple",
				Count:    1,
				MetaData: 0,
				Slot:     2,
			},
			"",
		)
		api.Replaceitem().ReplaceitemInContainerAsync(
			console.BlockPosByIndex(nbt_console.ConsoleIndexCenterBlock),
			game_interface.ReplaceitemInfo{
				Name:     "diamond_sword",
				Count:    1,
				MetaData: 200,
				Slot:     1,
			},
			"",
		)
		api.Commands().AwaitChangesGeneral()

		container := block_helper.ContainerBlockHelper{
			OpenInfo: block_helper.ContainerBlockOpenInfo{
				Name:                  "light_blue_shulker_box",
				States:                nil,
				ConsiderOpenDirection: true,
				ShulkerFacing:         0,
			},
			IsEmpty: false,
		}
		*console.BlockByIndex(nbt_console.ConsoleIndexCenterBlock) = container

		err := itemCache.StoreCache(
			[]item_cache.ItemCacheInfo{
				{SlotID: 2, Count: 1, Hash: item_cache.ItemHashNumber{HashNumber: 3, SetHashNumber: item_cache.SetHashNumberNotExist}},
				{SlotID: 1, Count: 1, Hash: item_cache.ItemHashNumber{HashNumber: 4, SetHashNumber: item_cache.SetHashNumberNotExist}},
			},
			container.OpenInfo,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 1 failed due to %v (stage 1)", err))
		}
	}

	// Test round 4
	{
		slotIDA, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    3,
				SetHashNumber: item_cache.SetHashNumberNotExist,
			},
			nil,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 4 failed due to %v (stage 1)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 4 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 4 failed")
		}

		api.Commands().SendWSCommandWithResp("give @s barrier 2048")
		for index := range 36 {
			console.UseInventorySlot(nbt_console.RequesterUser, resources_control.SlotID(index), true)
		}

		slotIDB, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    4,
				SetHashNumber: item_cache.SetHashNumberNotExist,
			},
			[]resources_control.SlotID{slotIDA},
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 4 failed due to %v (stage 2)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 4 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 4 failed")
		}
		if slotIDB != 1 {
			panic("SystemTestingItemCache: Test round 4 failed")
		}

		api.Commands().SendWSCommandWithResp("clear")
		console.CleanInventory()

		slotIDC, hit, isSetHashHit, err := itemCache.LoadCache(
			item_cache.ItemHashNumber{
				HashNumber:    3,
				SetHashNumber: item_cache.SetHashNumberNotExist,
			},
			nil,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemCache: Test round 4 failed due to %v (stage 3)", err))
		}
		if isSetHashHit {
			panic("SystemTestingItemCache: Test round 4 failed")
		}
		if !hit {
			panic("SystemTestingItemCache: Test round 4 failed")
		}
		if slotIDC != slotIDA {
			panic("SystemTestingItemCache: Test round 4 failed")
		}
	}

	pterm.Success.Printfln("SystemTestingItemCache: PASS (Time used = %v)", time.Since(tA))
}
