package nbt_parser_block

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_parser_general "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/general"
	"github.com/mitchellh/mapstructure"
)

type BannerNBT struct {
	Base     int32
	Patterns []nbt_parser_general.BannerPattern
	Type     int32
}

type Banner struct {
	DefaultBlock
	NBT BannerNBT
}

func (b Banner) NeedSpecialHandle() bool {
	if b.NBT.Base != nbt_parser_general.BannerBaseColorDefault {
		return true
	}
	if len(b.NBT.Patterns) > 0 {
		return true
	}
	if b.NBT.Type == nbt_parser_general.BannerTypeOminous {
		return true
	}
	return false
}

func (b Banner) NeedCheckCompletely() bool {
	return true
}

func (b *Banner) Parse(nbtMap map[string]any) error {
	patterns, _ := nbtMap["Patterns"].([]any)
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
			return fmt.Errorf("Parse: %v", err)
		}

		b.NBT.Patterns = append(b.NBT.Patterns, pattern)
	}

	b.NBT.Base, _ = nbtMap["Base"].(int32)
	b.NBT.Type, _ = nbtMap["Type"].(int32)

	return nil
}

func (b Banner) StableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)
	basicInfo := b.DefaultBlock.StableBytes()

	w.ByteSlice(&basicInfo)
	w.Varint32(&b.NBT.Base)
	protocol.SliceUint16Length(w, &b.NBT.Patterns)
	w.Varint32(&b.NBT.Type)

	return buf.Bytes()
}
