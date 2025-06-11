package nbt_block

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
)

type CommandBlock struct {
	NBTBlockBase
	data nbt_parser_block.CommandBlock
}

func (c *CommandBlock) Make() error {
	var mode uint32 = packet.CommandBlockImpulse
	api := c.console.API()

	if c.data.BlockName() == "minecraft:chain_command_block" {
		mode = packet.CommandBlockChain
	}
	if c.data.BlockName() == "minecraft:repeating_command_block" {
		mode = packet.CommandBlockRepeating
	}

	err := api.Resources().WritePacket(&packet.CommandBlockUpdate{
		Block:              true,
		Position:           c.console.Center(),
		Mode:               mode,
		NeedsRedstone:      c.data.NBT.Auto == 0,
		Conditional:        c.data.NBT.ConditionalMode == 1,
		Command:            c.data.NBT.Command,
		Name:               c.data.NBT.CustomName,
		ShouldTrackOutput:  c.data.NBT.TrackOutput == 1,
		TickDelay:          c.data.NBT.TickDelay,
		ExecuteOnFirstTick: c.data.NBT.ExecuteOnFirstTick == 1,
	})
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	err = api.Commands().AwaitChangesGeneral()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	return nil
}
