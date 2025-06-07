package main

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache/item_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/pterm/pterm"
)

func SystemTestingLogin() {
	var err error
	tA := time.Now()

	cfg := client.Config{
		AuthServerAddress:    "...",
		AuthServerToken:      "...",
		RentalServerCode:     "48285363",
		RentalServerPasscode: "",
	}

	c, err = client.LoginRentalServer(cfg)
	if err != nil {
		panic(err)
	}

	resources = resources_control.NewResourcesControl(c)
	api = game_interface.NewGameInterface(resources)

	console, err = nbt_console.NewConsole(api, [3]int32{23, 12, -21})
	if err != nil {
		panic(fmt.Sprintf("SystemTestingSetblock: Failed on init new console, and the err is %v", err))
	}
	itemCache = item_cache.NewItemCache(console)

	pterm.Success.Printfln("SystemTestingLogin: PASS (Time used = %v)", time.Since(tA))
}
