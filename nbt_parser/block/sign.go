package nbt_parser_block

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/mitchellh/mapstructure"
)

type SignText struct {
	IgnoreLighting byte   `mapstructure:"IgnoreLighting"`
	SignTextColor  int32  `mapstructure:"SignTextColor"`
	Text           string `mapstructure:"Text"`
}

type SignNBT struct {
	IsWaxed   byte     `mapstructure:"IsWaxed"`
	FrontText SignText `mapstructure:"FrontText"`
	BackText  SignText `mapstructure:"BackText"`
}

type Sign struct {
	DefaultBlock
	NBT SignNBT
}

func (s *Sign) NeedSpecialHandle() bool {
	if s.NBT.IsWaxed == 1 {
		return true
	}

	texts := []SignText{s.NBT.FrontText, s.NBT.BackText}
	for _, value := range texts {
		if len(value.Text) > 0 {
			return true
		}
		if value.SignTextColor != utils.EncodeVarRGBA(0, 0, 0, 255) {
			return true
		}
		if value.IgnoreLighting == 1 {
			return true
		}
	}

	return false
}

func (s Sign) NeedCheckCompletely() bool {
	return true
}

func (s *Sign) Parse(nbtMap map[string]any) error {
	var result SignNBT
	var legacy SignText

	if _, ok := nbtMap["IsWaxed"]; ok {
		err := mapstructure.Decode(&nbtMap, &result)
		if err != nil {
			return fmt.Errorf("(s *Sign) Parse: %v", err)
		}
		s.NBT = result
	} else {
		err := mapstructure.Decode(&nbtMap, &legacy)
		if err != nil {
			return fmt.Errorf("(s *Sign) Parse: %v", err)
		}
		s.NBT.FrontText = legacy
		s.NBT.BackText = SignText{
			IgnoreLighting: 0,
			SignTextColor:  utils.EncodeVarRGBA(0, 0, 0, 255),
			Text:           "",
		}
	}

	rgb, _ := utils.DecodeVarRGBA(s.NBT.FrontText.SignTextColor)
	bestColor := utils.SearchForBestColor(rgb, mapping.DefaultDyeColor)
	s.NBT.FrontText.SignTextColor = utils.EncodeVarRGBA(bestColor[0], bestColor[1], bestColor[2], 255)

	rgb, _ = utils.DecodeVarRGBA(s.NBT.BackText.SignTextColor)
	bestColor = utils.SearchForBestColor(rgb, mapping.DefaultDyeColor)
	s.NBT.BackText.SignTextColor = utils.EncodeVarRGBA(bestColor[0], bestColor[1], bestColor[2], 255)

	return nil
}

func (s Sign) StableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	basicInfo := s.DefaultBlock.StableBytes()
	w.ByteSlice(&basicInfo)
	w.Uint8(&s.NBT.IsWaxed)

	texts := []SignText{s.NBT.FrontText, s.NBT.BackText}
	for _, value := range texts {
		w.String(&value.Text)
		w.Int32(&value.SignTextColor)
		w.Uint8(&value.IgnoreLighting)
	}

	return buf.Bytes()
}
