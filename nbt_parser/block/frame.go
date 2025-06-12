package nbt_parser_block

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
)

// FrameNBT ..
type FrameNBT struct {
	ItemRotation float32
	HaveItem     bool
	Item         nbt_parser_interface.Item
}

// 物品展示框
type Frame struct {
	DefaultBlock
	NBT FrameNBT
}

func (f Frame) NeedSpecialHandle() bool {
	return f.NBT.HaveItem
}

func (f Frame) NeedCheckCompletely() bool {
	return true
}

func (f *Frame) Parse(nbtMap map[string]any) error {
	f.NBT.ItemRotation, _ = nbtMap["ItemRotation"].(float32)

	itemMap, ok := nbtMap["Item"].(map[string]any)
	if ok {
		item, err := nbt_parser_interface.ParseItemNormal(itemMap)
		if err != nil {
			return fmt.Errorf("Parse: %v", err)
		}
		if item.ItemName() != "minecraft:filled_map" {
			f.NBT.HaveItem = true
			f.NBT.Item = item
		}
	}

	return nil
}

func (f Frame) NBTStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	w.Bool(&f.NBT.HaveItem)
	if f.NBT.HaveItem {
		itemStableBytes := f.NBT.Item.TypeStableBytes()
		w.ByteSlice(&itemStableBytes)
	}

	return buf.Bytes()
}

func (f *Frame) FullStableBytes() []byte {
	return append(f.DefaultBlock.FullStableBytes(), f.NBTStableBytes()...)
}
