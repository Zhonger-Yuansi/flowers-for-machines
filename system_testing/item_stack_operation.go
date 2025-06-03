package main

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/pterm/pterm"
)

func SystemTestingItemStackOperation() {
	tA := time.Now()

	// Test round 1
	{
		api.Commands().SendSettingsCommand("clear", true)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s apple 25", true) // Slot 0 (0 -> 25)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s diamond_sword 1", true) // Slot 1 (0 -> 1)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s red_flower 20", true) // Slot 2 (0 -> 20)
		api.Commands().AwaitChangesGeneral()

		success, _, _, _ := api.ItemStackOperation().OpenTransaction().
			MoveBetweenInventory(0, 3, 24). // Slot 3 (0 -> 24), Slot 0 (25 -> 1)
			MoveBetweenInventory(0, 3, 1).  // Slot 0 -> Slot 3
			MoveBetweenInventory(1, 4, 1).  // Slot 1 -> Slot 4
			MoveBetweenInventory(3, 0, 25). // Slot 3 -> Slot 0
			MoveBetweenInventory(0, 3, 10). // Slot 3 (0 -> 10), Slot 0 (25 -> 15)
			MoveBetweenInventory(3, 5, 5).  // Slot 5 (0 -> 5); Slot 3 (10 -> 5)
			MoveToContainer(4, 0, 1).       // Slot 4 (Inventory) -> Slot 0 (Chest)
			MoveToContainer(5, 1, 5).       // Slot 5 (Inventory) -> Slot 1 (Chest)
			MoveToContainer(3, 1, 5).       // Slot 3 (Inventory) -> Slot 1 (Chest)
			MoveToContainer(2, 2, 20).      // Slot 2 (Inventory) -> Slot 2 (Chest)
			MoveToInventory(2, 8, 6).       // Slot 2 (Chest, 20 -> 14) -> Slot 8 (Inventory, 0 -> 6)
			DropInventoryItem(8, 3).        // Slot 8 (6 -> 3)
			MoveToContainer(8, 2, 3).       // Slot 8 (Inventory, 3 -> 0) -> Slot 2 (Chest, 14 -> 17)
			DropInventoryItem(0, 15).       // Slot 0 (15 -> 0)
			Commit()
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 1")
		}
	}

	// Test round 2
	{
		err := api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 2 failed due to %v (stage 1)", err))
		}

		api.Commands().SendSettingsCommand("clear", true)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s apple 25", true) // Slot 0 (0 -> 25)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s diamond_sword 1", true) // Slot 1 (0 -> 1)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s red_flower 20", true) // Slot 2 (0 -> 20)
		api.Commands().AwaitChangesGeneral()

		states, err := api.SetBlock().SetAnvil([3]int32{0, 0, 0}, true)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 2 failed due to %v (stage 2)", err))
		}

		success, err := api.ContainerOpenAndClose().OpenContainer(
			game_interface.UseItemOnBlocks{
				HotbarSlotID: 2,
				BlockPos:     [3]int32{0, 0, 0},
				BlockName:    "anvil",
				BlockStates:  states,
			},
			true,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 2 failed due to %v (stage 3)", err))
		}
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 2")
		}

		success, _, _, _ = api.ItemStackOperation().OpenTransaction().
			MoveBetweenInventory(0, 3, 25).                 // Slot 0 -> Slot 3
			MoveBetweenInventory(1, 4, 1).                  // Slot 1 -> Slot 4
			MoveBetweenInventory(2, 5, 20).                 // Slot 2 -> Slot 5
			RenameInventoryItem(3, 25, "SYSTEM TESTING A"). // Hacking Attempt
			RenameInventoryItem(4, 1, "SYSTEM TESTING B").  // Hacking Attempt
			SwapBetweenInventory(3, 4).                     // Slot 3 <-> Slot 4
			SwapBetweenInventory(3, 5).                     // Slot 3 <-> Slot 5
			RenameInventoryItem(5, 1, "INLINE").            // Hacking Attempt
			RenameInventoryItem(5, 1, "INLINE A").          // Hacking Attempt
			RenameInventoryItem(3, 20, "INLINE B").         // Hacking Attempt
			RenameInventoryItem(4, 25, "INLINE C").         // Hacking Attempt
			RenameInventoryItem(3, 20, "APPLE").            // Real Name
			RenameInventoryItem(4, 25, "SWORD").            // Real Name
			RenameInventoryItem(5, 1, "献给机械の花束").           // Real Name
			SwapBetweenInventory(3, 4).                     // Slot 3 <-> Slot 4
			SwapBetweenInventory(4, 5).                     // Slot 4 <-> Slot 5
			MoveToContainer(3, 1, 25).                      // Slot 3 (Inventory) -> Slot 1 (Anvil)
			DropContainerItem(1, 25).                       // Slot 1 (Anvil, 25 -> 0)
			SwapBetweenInventory(4, 5).                     // Slot 4 <-> Slot 5
			SwapBetweenInventory(4, 5).                     // Slot 4 <-> Slot 5
			SwapBetweenInventory(5, 4).                     // Slot 5 <-> Slot 4
			DropInventoryItem(5, 1).                        // Slot 5 (1 -> 0)
			MoveBetweenInventory(4, 5, 10).                 // Slot 4 (20 -> 10) -> Slot 5 (0 -> 10)
			MoveToContainer(4, 1, 10).                      // Slot 4 (Inventory) -> Slot 1 (Anvil)
			SwapInventoryBetweenContainer(5, 1).            // Slot 5 (Anvil) <-> Slot 1 (Anvil)
			MoveToInventory(1, 4, 10).                      // Slot 1 (Anvil) -> Slot 4 (Inventory)
			DropInventoryItem(4, 10).                       // Slot 4 (10 -> 0)
			DropInventoryItem(5, 10).                       // Slot 5 (10 -> 0)
			Commit()
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 2")
		}

		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 2 failed due to %v (stage 4)", err))
		}
	}

	// Test round 3
	{
		success, err := api.ContainerOpenAndClose().OpenInventory()
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 3 failed due to %v (stage 1)", err))
		}
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 3")
		}

		api.Commands().SendSettingsCommand("clear", true)
		api.Commands().AwaitChangesGeneral()
		success, _, _, _ = api.ItemStackOperation().OpenTransaction().
			GetCreativeItemToInventory(1, 0, 64).
			GetCreativeItemToInventory(2, 1, 64).
			GetCreativeItemToInventory(0x5bc, 8, 1).
			DropInventoryItem(0, 64).
			DropInventoryItem(1, 64).
			DropInventoryItem(8, 1).
			Commit()
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 3")
		}

		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 3 failed due to %v (stage 2)", err))
		}
	}

	// Test round 4
	{
		api.Commands().SendSettingsCommand("clear", true)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s banner 1 10", true) // Slot 0
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s yellow_dye 20", true) // Slot 1
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s mojang_banner_pattern 1", true) // Slot 2
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s red_dye 20", true) // Slot 3
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s skull_banner_pattern", true) // Slot 4
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s bordure_indented_banner_pattern", true) // Slot 5
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s banner 1 11", true) // Slot 6
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand("give @s light_blue_dye 20", true) // Slot 7
		api.Commands().AwaitChangesGeneral()

		err := api.SetBlock().SetBlock(protocol.BlockPos{0, 0, 0}, "loom", `["direction"=0]`)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 4 failed due to %v (stage 1)", err))
		}

		success, err := api.ContainerOpenAndClose().OpenContainer(
			game_interface.UseItemOnBlocks{
				HotbarSlotID: 2,
				BlockPos:     protocol.BlockPos{0, 0, 0},
				BlockName:    "loom",
				BlockStates: map[string]any{
					"direction": int32(0),
				},
			},
			false,
		)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 4 failed due to %v (stage 2)", err))
		}
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 4")
		}

		success, _, _, _ = api.ItemStackOperation().OpenTransaction().
			LoomingFromInventory("bo", 0, 0, 1, resources_control.ExpectedNewItem{NetworkID: -1}).  // Banner 1 (1)
			LoomingFromInventory("moj", 2, 0, 3, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 1 (2)
			LoomingFromInventory("sku", 4, 0, 1, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 1 (3)
			LoomingFromInventory("sku", 4, 0, 1, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 1 (4)
			LoomingFromInventory("sku", 4, 0, 1, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 1 (5)
			LoomingFromInventory("sku", 4, 0, 1, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 1 (6)
			LoomingFromInventory("cbo", 5, 6, 3, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 2 (1)
			LoomingFromInventory("", 0, 6, 1, resources_control.ExpectedNewItem{NetworkID: -1}).    // Banner 2 (2)
			LoomingFromInventory("moj", 2, 6, 7, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 2 (3)
			LoomingFromInventory("sku", 4, 6, 3, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 2 (4)
			LoomingFromInventory("cbo", 5, 6, 1, resources_control.ExpectedNewItem{NetworkID: -1}). // Banner 2 (5)
			LoomingFromInventory("bo", 0, 6, 1, resources_control.ExpectedNewItem{NetworkID: -1}).  // Banner 2 (6)
			DropInventoryItem(0, 1).
			DropInventoryItem(6, 1).
			Commit()
		if !success {
			panic("SystemTestingItemStackOperation: Failed on test round 4")
		}

		err = api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			panic(fmt.Sprintf("SystemTestingItemStackOperation: Test round 4 failed due to %v (stage 3)", err))
		}
	}

	pterm.Success.Printfln("SystemTestingItemStackOperation: PASS (Time used = %v)", time.Since(tA))
}
