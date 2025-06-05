package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/pterm/pterm"
)

func SystemTestingBotClick() {
	tA := time.Now()

	// ClickBlock
	{
		api.Commands().SendSettingsCommand("gamemode 0", true)
		api.Commands().SendSettingsCommand("tp 0 0 0", true)
		api.Commands().SendSettingsCommand("clear", true)
		api.Commands().AwaitChangesGeneral()

		api.Commands().SendSettingsCommand("replaceitem entity @s slot.hotbar 2 apple 10", true)
		api.BotClick().ChangeSelectedHotbarSlot(2)
		api.Commands().SendSettingsCommand("setblock 0 0 0 air", true)
		api.Commands().SendSettingsCommand("setblock 0 -1 0 grass", true)
		api.Commands().AwaitChangesGeneral()
		api.Commands().SendSettingsCommand(`setblock 0 0 0 glow_frame ["facing_direction"=1]`, true)
		api.Commands().AwaitChangesGeneral()

		channel := make(chan struct{})
		uniqueID := api.Resources().Inventories().SetCallback(
			0, 2,
			func(item *protocol.ItemInstance) {
				if item.Stack.Count != 9 {
					panic("SystemTestingBotClick: `ClickBlock` failed")
				}
				close(channel)
			},
		)

		api.BotClick().ClickBlock(
			game_interface.UseItemOnBlocks{
				HotbarSlotID: 2,
				BlockPos:     [3]int32{0, 0, 0},
				BlockName:    "glow_frame",
				BlockStates: map[string]any{
					"facing_direction":     int32(1),
					"item_frame_map_bit":   byte(0),
					"item_frame_photo_bit": byte(0),
				},
			},
		)

		timer := time.NewTimer(time.Second * 5)
		defer timer.Stop()
		select {
		case <-timer.C:
			panic("SystemTestingBotClick: `ClickBlock` time out")
		case <-channel:
			api.PacketListener().DestroyListener(uniqueID)
		}
	}

	// PickBlock
	{
		api.Commands().SendSettingsCommand("gamemode 1", true)
		api.Commands().SendSettingsCommand("clear", true)
		api.Commands().AwaitChangesGeneral()

		success, resultHotbar, err := api.BotClick().PickBlock([3]int32{0, 0, 0}, true)
		if err != nil {
			panic(fmt.Sprintf("SystemTestingBotClick: `PickBlock` failed due to %v", err))
		}
		if !success {
			panic("SystemTestingBotClick: `PickBlock` failed on test round 1")
		}
		if resultHotbar != 0 {
			panic("SystemTestingBotClick: `PickBlock` failed on test round 2")
		}

		item, _ := api.Resources().Inventories().GetItemStack(0, 0)
		if item == nil {
			panic("SystemTestingBotClick: `PickBlock` failed on test round 3")
		}
		if !strings.Contains(fmt.Sprintf("%#v", item.Stack.NBTData), "(+DATA)") {
			panic("SystemTestingBotClick: `PickBlock` failed on test round 4")
		}
	}

	pterm.Success.Printfln("SystemTestingBotClick: PASS (Time used = %v)", time.Since(tA))
}
