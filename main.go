package main

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

func main() {
	cfg := client.Config{
		AuthServerAddress:    "...",
		AuthServerToken:      "...",
		RentalServerCode:     "48285363",
		RentalServerPasscode: "",
	}

	c, err := client.LoginRentalServer(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		c.Conn().Close()
		time.Sleep(time.Second)
	}()

	resources := resources_control.NewResourcesControl(c)
	api := game_interface.NewGameInterface(resources)

	api.Commands().SendSettingsCommand("clear", true)
	api.Commands().SendSettingsCommand("give @s apple 10", true)
	api.Commands().SendSettingsCommand("give @s diamond_sword 1", true)
	api.Commands().SendSettingsCommand(`setblock 0 0 0 anvil ["minecraft:cardinal_direction"="north","damage"="undamaged"]`, true)
	// api.Commands().SendSettingsCommand("tp 0 0 0", true)
	api.BotClick().ChangeSelectedHotbarSlot(0)
	api.Commands().AwaitChangesGeneral()

	channel := make(chan struct{})
	api.Resources().Container().SetContainerOpenCallback(func() { close(channel) })
	api.BotClick().ClickBlock(game_interface.UseItemOnBlocks{
		HotbarSlotID: 0,
		BlockPos:     [3]int32{0, 0, 0},
		BlockName:    "anvil",
		BlockStates: map[string]any{
			`minecraft:cardinal_direction`: "north",
			"damage":                       "undamaged",
		},
	})
	<-channel

	tran := api.ItemStackOperation().OpenTransaction()
	tran.RenameInventoryItem(1, "haha Testing")
	tran.RenameInventoryItem(0, "System Testing")
	tran.MoveInventoryItem(0, 2, 1)
	tran.MoveInventoryItem(1, 3, 1)
	fmt.Println(tran.Commit())
}

// func legacy() {
// 	// {
// 	// 	api.BotClick().ChangeSelectedHotbarSlot(0)
// 	// 	api.Commands().SendWSCommandWithResp("tp 0 0 0")
// 	// 	api.BotClick().ClickBlock(game_interface.UseItemOnBlocks{
// 	// 		HotbarSlotID: 0,
// 	// 		BlockPos:     [3]int32{0, 0, 0},
// 	// 		BlockName:    "dispenser",
// 	// 		BlockStates: map[string]any{
// 	// 			`facing_direction`: int32(0),
// 	// 			"triggered_bit":    false,
// 	// 		},
// 	// 	})
// 	// 	api.Commands().AwaitChangesGeneral()
// 	// }

// 	{
// 		api.Commands().SendSettingsCommand("clear", true)
// 		api.Commands().SendSettingsCommand("gamemode 1", true)
// 		api.Commands().SendSettingsCommand("give @s apple 10", true)
// 		api.Commands().AwaitChangesGeneral()

// 		// channel := make(chan struct{})
// 		// api.Resources().Container().SetContainerOpenCallback(func() { close(channel) })
// 		// api.Resources().WritePacket(&packet.Interact{
// 		// 	ActionType:            packet.InteractActionOpenInventory,
// 		// 	TargetEntityRuntimeID: api.GetBotInfo().EntityRuntimeID,
// 		// })
// 		// <-channel

// 		i, _ := api.Resources().Inventories().GetItemStack(0, 0)

// 		t1 := protocol.PlaceStackRequestAction{}
// 		t1.Count = 10
// 		t1.Source = protocol.StackRequestSlotInfo{
// 			ContainerID:    protocol.ContainerHotBar,
// 			Slot:           0,
// 			StackNetworkID: i.StackNetworkID,
// 		}
// 		t1.Destination = protocol.StackRequestSlotInfo{
// 			ContainerID:    protocol.ContainerHotBar,
// 			Slot:           1,
// 			StackNetworkID: 0,
// 		}

// 		t2 := protocol.PlaceStackRequestAction{}
// 		t2.Count = 10
// 		t2.Source = protocol.StackRequestSlotInfo{
// 			ContainerID:    protocol.ContainerHotBar,
// 			Slot:           1,
// 			StackNetworkID: -1,
// 		}
// 		t2.Destination = protocol.StackRequestSlotInfo{
// 			ContainerID:    protocol.ContainerHotBar,
// 			Slot:           2,
// 			StackNetworkID: 0,
// 		}

// 		api.Resources().WritePacket(&packet.ItemStackRequest{
// 			Requests: []protocol.ItemStackRequest{
// 				{
// 					RequestID: -1,
// 					Actions: []protocol.StackRequestAction{
// 						&t1,
// 					},
// 				},
// 				{
// 					RequestID: -3,
// 					Actions: []protocol.StackRequestAction{
// 						&t2,
// 					},
// 				},
// 			},
// 		})
// 		api.Commands().AwaitChangesGeneral()
// 		return
// 	}

// 	err = api.Commands().SendPlayerCommand("tp 0 0 0")
// 	fmt.Println(err)

// 	resp, err := api.Commands().SendWSCommandWithResp("say 123")
// 	fmt.Println(resp, err)

// 	resp, isTimeout, err := api.Commands().SendPlayerCommandWithTimeout("say 123", time.Second*5)
// 	fmt.Println(resp, isTimeout, err)

// 	querytargetResult, err := api.Querytarget().DoQuerytarget("@s")
// 	fmt.Println(querytargetResult, err)

// 	err = api.Commands().SendSettingsCommand(`setblock 0 0 0 anvil ["minecraft:cardinal_direction"="north","damage"="undamaged"]`, true)
// 	fmt.Println(err)
// 	err = api.Commands().AwaitChangesGeneral()
// 	fmt.Println(err)

// 	{
// 		channel := make(chan struct{})
// 		_ = api.Resources().PacketListener().ListenPacket(
// 			[]uint32{packet.IDContainerOpen},
// 			func(p packet.Packet) bool {
// 				close(channel)
// 				fmt.Println("Container opened")
// 				return true
// 			},
// 		)

// 		err = api.BotClick().ChangeSelectedHotbarSlot(0)
// 		fmt.Println(err)
// 		err = api.BotClick().ClickBlock(game_interface.UseItemOnBlocks{
// 			HotbarSlotID: 0,
// 			BlockPos:     [3]int32{0, 0, 0},
// 			BlockName:    "anvil",
// 			BlockStates: map[string]any{
// 				"minecraft:cardinal_direction": "north",
// 				"damage":                       "undamaged",
// 			},
// 		})
// 		fmt.Println(err)

// 		<-channel
// 		fmt.Println(api.Resources().Container().ContainerData())
// 	}

// 	{
// 		channel := make(chan struct{})
// 		api.Resources().Container().SetContainerCloseCallback(
// 			func(isServerSide bool) {
// 				close(channel)
// 				fmt.Println("Container closed")
// 			},
// 		)

// 		containerData, existed := api.Resources().Container().ContainerData()
// 		fmt.Println(existed)
// 		err = api.Resources().WritePacket(&packet.ContainerClose{
// 			WindowID: containerData.WindowID,
// 		})
// 		fmt.Println(err)

// 		<-channel
// 		fmt.Println(api.Resources().Container().ContainerData())
// 	}

// 	{
// 		uniqueID, err := api.StructureBackup().BackupStructure([3]int32{0, 0, 0})
// 		fmt.Println(uniqueID, err)

// 		err = api.StructureBackup().RevertStructure(uniqueID, [3]int32{0, 1, 0})
// 		fmt.Println(err)
// 		err = api.StructureBackup().DeleteStructure(uniqueID)
// 		fmt.Println(err)
// 	}

// 	{
// 		_, err = api.Commands().SendPlayerCommandWithResp("tp 27 -60 -79")
// 		fmt.Println(err)
// 		fmt.Println(api.BotClick().PickBlock([3]int32{27, -60, -79}, 0, true))
// 	}

// 	{
// 		err = api.Commands().SendSettingsCommand("replaceitem entity @s slot.hotbar 1 banner 1 10", true)
// 		fmt.Println(err)

// 		err = api.BotClick().ChangeSelectedHotbarSlot(1)
// 		fmt.Println(err)
// 		err = api.Commands().AwaitChangesGeneral()
// 		fmt.Println(err)

// 		_, err = api.BotClick().PlaceBlockHighLevel(
// 			[3]int32{24, -57, -75},
// 			1,
// 			2,
// 		)
// 		fmt.Println(err)
// 	}

// 	api.Commands().SendChat("aaaa")
// }
