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

	resp, err := api.Commands().SendWSCommandWithResp("say 123")
	fmt.Println(resp, err)

	resp, isTimeout, err := api.Commands().SendPlayerCommandWithTimeout("say 123", time.Second*5)
	fmt.Println(resp, isTimeout, err)

	uniqueID, err := api.StructureBackup().BackupStructure([3]int32{0, 0, 0})
	fmt.Println(uniqueID, err)

	err = api.StructureBackup().RevertStructure(uniqueID, [3]int32{0, 1, 0})
	fmt.Println(err)
	err = api.StructureBackup().DeleteStructure(uniqueID)
	fmt.Println(err)

	resp, err = api.Commands().SendWSCommandWithResp("querytarget @s")
	fmt.Println(resp, err)

	querytargetResult, err := api.Querytarget().DoQuerytarget("@s")
	fmt.Println(querytargetResult, err)

	api.Commands().SendChat("aaaa")
}
