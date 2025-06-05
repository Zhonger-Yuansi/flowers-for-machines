package main

import (
	"time"

	"github.com/pterm/pterm"
)

func SystemTestingConsole() {
	tA := time.Now()

	// Test round 1
	{
		err := console.InitConsoleArea()
		if err != nil {
			panic("SystemTestingSetblock: Test round 1 failed")
		}
	}

	pterm.Success.Printfln("SystemTestingConsole: PASS (Time used = %v)", time.Since(tA))
}
