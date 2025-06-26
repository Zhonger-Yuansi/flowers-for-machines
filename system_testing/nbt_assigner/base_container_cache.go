package main

import (
	"fmt"
	"time"

	"github.com/OmineDev/flowers-for-machines/nbt_assigner/block_helper"
	"github.com/OmineDev/flowers-for-machines/nbt_assigner/nbt_console"
	"github.com/OmineDev/flowers-for-machines/utils"

	"github.com/pterm/pterm"
)

func SystemTestingBaseContainerCache() {
	tA := time.Now()
	barrelStatesString := `["facing_direction"=1,"open_bit"=false]`

	// Test round 1
	{
		api.SetBlock().SetBlock(console.Center(), "barrel", barrelStatesString)

		container := block_helper.ContainerBlockHelper{
			OpenInfo: block_helper.ContainerBlockOpenInfo{
				Name:                  "barrel",
				States:                utils.ParseBlockStatesString(barrelStatesString),
				ConsiderOpenDirection: false,
			},
		}
		console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, container)

		err := baseContainerCache.StoreCache("", 0)
		if err != nil {
			panic("SystemTestingBaseContainerCache: Failed on test round 1")
		}
	}

	// Test round 2
	{
		api.SetBlock().SetBlock(console.Center(), "air", "[]")
		console.UseHelperBlock(nbt_console.RequesterUser, nbt_console.ConsoleIndexCenterBlock, block_helper.Air{})

		hit, _ := baseContainerCache.LoadCache("barrel", utils.ParseBlockStatesString(barrelStatesString), "", 0)
		if !hit {
			panic("SystemTestingBaseContainerCache: Failed on test round 2")
		}

		success, _ := console.OpenContainerByIndex(nbt_console.ConsoleIndexCenterBlock)
		if !success {
			panic("SystemTestingBaseContainerCache: Failed on test round 2")
		}

		err := api.ContainerOpenAndClose().CloseContainer()
		if err != nil {
			panic(fmt.Sprintf("SystemTestingBaseContainerCache: Failed on test round 2 due to %v", err))
		}
	}

	pterm.Success.Printfln("SystemTestingBaseContainerCache: PASS (Time used = %v)", time.Since(tA))
}
