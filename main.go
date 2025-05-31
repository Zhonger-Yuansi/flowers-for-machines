package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
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

	for {
		pk := <-c.CachedPacket()
		if pk == nil {
			break
		}
		fmt.Println(pk.ID())
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for {
			pk, err := c.Conn().ReadPacket()
			if err != nil {
				panic(err)
			}
			if p, ok := pk.(*packet.CommandOutput); ok {
				fmt.Printf("%#v\n", p)
				wg.Done()
				return
			} else {
				pterm.Info.Println(pk)
			}
		}
	}()

	time.Sleep(time.Second)
	fmt.Println("SEND")

	c.Conn().WritePacket(&packet.CommandRequest{
		CommandLine: "System Testing",
		CommandOrigin: protocol.CommandOrigin{
			Origin:    protocol.CommandOriginAutomationPlayer,
			UUID:      uuid.New(),
			RequestID: "96045347-a6a3-4114-94c0-1bc4cc561694",
		},
		Internal:  false,
		UnLimited: false,
		Version:   0x24,
	})

	wg.Wait()
}
