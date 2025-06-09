package nbt_parser_block

import (
	"bytes"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// 默认 NBT 实体
type DefaultBlock struct {
	Name   string
	States map[string]any
}

func (d *DefaultBlock) BlockName() string {
	d.Name = strings.ToLower(d.Name)
	if !strings.HasPrefix(d.Name, "minecraft:") {
		d.Name = "minecraft:" + d.Name
	}
	return d.Name
}

func (d DefaultBlock) BlockStates() map[string]any {
	return d.States
}

func (d DefaultBlock) BlockStatesString() string {
	return utils.MarshalBlockStates(d.States)
}

func (*DefaultBlock) Parse(nbtMap map[string]any) error {
	return nil
}

func (DefaultBlock) NeedSpecialHandle() bool {
	return false
}

func (DefaultBlock) NeedCheckCompletely() bool {
	return false
}

func (d DefaultBlock) StableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	name := d.BlockName()
	states := d.BlockStatesString()
	w.String(&name)
	w.String(&states)

	return buf.Bytes()
}
