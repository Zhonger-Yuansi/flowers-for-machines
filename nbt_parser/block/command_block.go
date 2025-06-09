package nbt_parser_block

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/mitchellh/mapstructure"
)

type CommandBlockNBT struct {
	Command            string `mapstructure:"Command"`
	CustomName         string `mapstructure:"CustomName"`
	TickDelay          int32  `mapstructure:"TickDelay"`
	ExecuteOnFirstTick byte   `mapstructure:"ExecuteOnFirstTick"`
	TrackOutput        byte   `mapstructure:"TrackOutput"`
	ConditionalMode    byte   `mapstructure:"conditionalMode"`
	Auto               byte   `mapstructure:"auto"`
	Version            int32  `mapstructure:"Version"`
}

type CommandBlock struct {
	DefaultBlock
	NBT CommandBlockNBT
}

func (c *CommandBlock) NeedSpecialHandle() bool {
	if len(c.NBT.Command) > 0 || len(c.NBT.CustomName) > 0 {
		return true
	}
	if c.NBT.TickDelay != 0 {
		return true
	}

	switch c.BlockName() {
	case "minecraft:repeating_command_block":
		if c.NBT.ExecuteOnFirstTick == 0 {
			return true
		}
	default:
		if c.NBT.ExecuteOnFirstTick == 1 {
			return true
		}
	}

	switch c.BlockName() {
	case "minecraft:chain_command_block":
		if c.NBT.Auto == 0 {
			return true
		}
	default:
		if c.NBT.Auto == 1 {
			return true
		}
	}

	return false
}

func (c CommandBlock) NeedCheckCompletely() bool {
	return false
}

func (c *CommandBlock) Parse(nbtMap map[string]any) error {
	var result CommandBlockNBT
	err := mapstructure.Decode(&nbtMap, &result)
	if err != nil {
		return fmt.Errorf("Parse: %v", err)
	}
	c.NBT = result
	return nil
}

func (c CommandBlock) StableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)
	basicInfo := c.DefaultBlock.StableBytes()

	w.ByteSlice(&basicInfo)
	w.String(&c.NBT.Command)
	w.String(&c.NBT.CustomName)
	w.Varint32(&c.NBT.TickDelay)
	w.Uint8(&c.NBT.ExecuteOnFirstTick)
	w.Uint8(&c.NBT.Auto)

	return buf.Bytes()
}
