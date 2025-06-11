package nbt_parser_item

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_general "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/general"
	"github.com/mitchellh/mapstructure"
)

// BannerNBT ..
type BannerNBT struct {
	Patterns []nbt_parser_general.BannerPattern
	Type     int32
}

// 旗帜
type Banner struct {
	DefaultItem
	NBT BannerNBT
}

// parse ..
func (b *Banner) parse(tag map[string]any) error {
	b.DefaultItem.Enhance.ItemComponent.LockInInventory = false
	b.DefaultItem.Enhance.ItemComponent.LockInSlot = false
	b.DefaultItem.Enhance.EnchList = nil
	b.DefaultItem.Block = ItemBlockData{}

	if len(tag) == 0 {
		return nil
	}

	patterns, _ := tag["Patterns"].([]any)
	if len(patterns) > 6 {
		patterns = patterns[0:6]
	}

	for _, value := range patterns {
		var pattern nbt_parser_general.BannerPattern

		val, ok := value.(map[string]any)
		if !ok {
			continue
		}

		err := mapstructure.Decode(&val, &pattern)
		if err != nil {
			return fmt.Errorf("parse: %v", err)
		}

		if mapping.BannerPatternUnsupported[pattern.Pattern] {
			continue
		}

		b.NBT.Patterns = append(b.NBT.Patterns, pattern)
	}

	b.NBT.Type, _ = tag["Type"].(int32)
	if b.NBT.Type == nbt_parser_general.BannerTypeOminous {
		b.NBT.Patterns = nil
	}

	return nil
}

func (b *Banner) ParseNormal(nbtMap map[string]any) error {
	tag, _ := nbtMap["tag"].(map[string]any)
	err := b.parse(tag)
	if err != nil {
		return fmt.Errorf("ParseNormal: %v", err)
	}
	return nil
}

func (b *Banner) ParseNetwork(item protocol.ItemStack, itemName string) error {
	err := b.parse(item.NBTData)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	return nil
}

func (b Banner) IsComplex() bool {
	if len(b.NBT.Patterns) > 0 {
		return true
	}
	if b.NBT.Type == nbt_parser_general.BannerTypeOminous {
		return true
	}
	return false
}

func (Banner) NeedCheckCompletely() bool {
	return false
}

func (b Banner) NBTStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	protocol.SliceUint16Length(w, &b.NBT.Patterns)
	w.Varint32(&b.NBT.Type)

	return buf.Bytes()
}

func (b *Banner) TypeStableBytes() []byte {
	return append(b.DefaultItem.TypeStableBytes(), b.NBTStableBytes()...)
}

func (b *Banner) FullStableBytes() []byte {
	return append(b.TypeStableBytes(), b.Basic.Count)
}
