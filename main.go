package main

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
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

	err = api.Commands().SendPlayerCommand("tp 0 0 0")
	fmt.Println(err)

	resp, err := api.Commands().SendWSCommandWithResp("say 123")
	fmt.Println(resp, err)

	resp, isTimeout, err := api.Commands().SendPlayerCommandWithTimeout("say 123", time.Second*5)
	fmt.Println(resp, isTimeout, err)

	querytargetResult, err := api.Querytarget().DoQuerytarget("@s")
	fmt.Println(querytargetResult, err)

	err = api.Commands().SendSettingsCommand(`setblock 0 0 0 anvil ["minecraft:cardinal_direction"="north","damage"="undamaged"]`, true)
	fmt.Println(err)
	err = api.Commands().AwaitChangesGeneral()
	fmt.Println(err)

	{
		channel := make(chan struct{})
		_ = api.Resources().PacketListener().ListenPacket(
			[]uint32{packet.IDContainerOpen},
			func(p packet.Packet) bool {
				close(channel)
				fmt.Println("Container opened")
				return true
			},
		)

		err = api.BotClick().ChangeSelectedHotbarSlot(0)
		fmt.Println(err)
		err = api.BotClick().ClickBlock(game_interface.UseItemOnBlocks{
			HotbarSlotID: 0,
			BlockPos:     [3]int32{0, 0, 0},
			BlockName:    "anvil",
			BlockStates: map[string]any{
				"minecraft:cardinal_direction": "north",
				"damage":                       "undamaged",
			},
		})
		fmt.Println(err)

		<-channel
		fmt.Println(api.Resources().Container().ContainerData())
	}

	{
		channel := make(chan struct{})
		api.Resources().Container().SetContainerCloseCallback(
			func(isServerSide bool) {
				close(channel)
				fmt.Println("Container closed")
			},
		)

		containerData, existed := api.Resources().Container().ContainerData()
		fmt.Println(existed)
		err = api.Resources().WritePacket(&packet.ContainerClose{
			WindowID: containerData.WindowID,
		})
		fmt.Println(err)

		<-channel
		fmt.Println(api.Resources().Container().ContainerData())
	}

	{
		uniqueID, err := api.StructureBackup().BackupStructure([3]int32{0, 0, 0})
		fmt.Println(uniqueID, err)

		err = api.StructureBackup().RevertStructure(uniqueID, [3]int32{0, 1, 0})
		fmt.Println(err)
		err = api.StructureBackup().DeleteStructure(uniqueID)
		fmt.Println(err)
	}

	{
		_, err = api.Commands().SendPlayerCommandWithResp("tp 27 -60 -79")
		fmt.Println(err)
		fmt.Println(api.BotClick().PickBlock([3]int32{27, -60, -79}, 0, true))
	}

	api.Commands().SendChat("aaaa")
}
