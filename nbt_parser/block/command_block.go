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
	LastOutput         string `mapstructure:"LastOutput"`
	TickDelay          int32  `mapstructure:"TickDelay"`
	ExecuteOnFirstTick bool   `mapstructure:"ExecuteOnFirstTick"`
	TrackOutput        bool   `mapstructure:"TrackOutput"`
	ConditionalMode    bool   `mapstructure:"conditionalMode"`
	Auto               bool   `mapstructure:"auto"`
	Version            int32  `mapstructure:"Version"`
}

type CommandBlock struct {
	DefaultBlock
	NBT CommandBlockNBT
}

func (c CommandBlock) NeedSpecialHandle() bool {
	return true
}

func (c CommandBlock) NeedCheckCompletely() bool {
	return false
}

func (c *CommandBlock) Parse(nbtMap map[string]any) error {
	var result CommandBlock
	err := mapstructure.Decode(&nbtMap, &result)
	if err != nil {
		return fmt.Errorf("(c *CommandBlock) Parse: %v", err)
	}
	*c = result
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
	w.Bool(&c.NBT.ExecuteOnFirstTick)
	w.Bool(&c.NBT.TrackOutput)
	w.Bool(&c.NBT.ConditionalMode)
	w.Bool(&c.NBT.Auto)

	return buf.Bytes()
}
