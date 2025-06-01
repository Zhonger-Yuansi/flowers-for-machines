package main

import (
	"fmt"
	"time"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/google/uuid"
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

	requestID := uuid.New()

	var resp *packet.CommandOutput
	channel := make(chan struct{})

	resources.Commands().SetCommandRequestCallback(
		requestID,
		func(p *packet.CommandOutput) {
			resp = p
			close(channel)
		},
	)

	resources.WritePacket()(&packet.CommandRequest{
		CommandLine: "System Testing",
		CommandOrigin: protocol.CommandOrigin{
			Origin:    protocol.CommandOriginAutomationPlayer,
			UUID:      requestID,
			RequestID: "96045347-a6a3-4114-94c0-1bc4cc561694",
		},
		Internal:  false,
		UnLimited: false,
		Version:   0x24,
	})
	<-channel
	fmt.Printf("%#v\n", resp)
}
