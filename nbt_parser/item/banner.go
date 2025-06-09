package nbt_parser

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/mitchellh/mapstructure"
)

// 描述旗帜的种类
const (
	BannerTypeNormal  int32 = iota // 普通旗帜
	BannerTypeOminous              // 不祥旗帜
)

// BannerPattern 是旗帜的单个图案
type BannerPattern struct {
	Color   int32  `mapstructure:"Color"`
	Pattern string `mapstructure:"Pattern"`
}

// Marshal ..
func (b *BannerPattern) Marshal(io protocol.IO) {
	io.Varint32(&b.Color)
	io.String(&b.Pattern)
}

type BannerNBT struct {
	Patterns []BannerPattern
	Type     int32
}

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
	for _, value := range patterns {
		var pattern BannerPattern

		val, ok := value.(map[string]any)
		if !ok {
			continue
		}

		err := mapstructure.Decode(&val, &pattern)
		if err != nil {
			return fmt.Errorf("parse: %v", err)
		}

		b.NBT.Patterns = append(b.NBT.Patterns, pattern)
	}

	b.NBT.Type, _ = tag["Type"].(int32)
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

func (b *Banner) ParseNetwork(item protocol.ItemStack, itemNetworkIDToName map[int32]string) error {
	err := b.parse(item.NBTData)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	return nil
}

func (b Banner) NeedSpecialHandle() bool {
	if b.DefaultItem.NeedSpecialHandle() {
		return true
	}
	if len(b.NBT.Patterns) > 0 {
		return true
	}
	if b.NBT.Type == BannerTypeOminous {
		return true
	}
	return false
}

func (Banner) NeedCheckCompletely() bool {
	return false
}

func (b *Banner) TypeStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	basicInfo := b.DefaultItem.TypeStableBytes()
	w.ByteSlice(&basicInfo)
	protocol.SliceUint16Length(w, &b.NBT.Patterns)
	w.Varint32(&b.NBT.Type)

	return buf.Bytes()
}

func (b *Banner) FullStableBytes() []byte {
	return append(b.TypeStableBytes(), b.Basic.Count)
}
