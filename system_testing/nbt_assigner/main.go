package main

import (
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache/base_container_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/pterm/pterm"
)

var (
	c         *client.Client
	resources *resources_control.Resources
	api       *game_interface.GameInterface
)

var (
	console            *nbt_console.Console
	baseContainerCache *base_container_cache.BaseContainerCache
)

func main() {
	tA := time.Now()

	SystemTestingLogin()
	defer func() {
		c.Conn().Close()
		time.Sleep(time.Second)
	}()

	SystemTestingBaseContainerCache()

	pterm.Success.Printfln("System Testing: ALL PASS (Time used = %v)", time.Since(tA))
}
