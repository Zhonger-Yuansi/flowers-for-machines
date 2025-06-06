package main

import (
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	"github.com/pterm/pterm"
)

func SystemTestingConsole() {
	var err error
	tA := time.Now()

	// Prepare
	console, err = nbt_console.NewConsole(api, [3]int32{23, 12, -21})
	if err != nil {
		panic("SystemTestingSetblock: Test round 1 failed")
	}

	pterm.Success.Printfln("SystemTestingConsole: PASS (Time used = %v)", time.Since(tA))
}
