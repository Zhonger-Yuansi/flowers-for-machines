package nbt_parser_block

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/mitchellh/mapstructure"
)

// StructureBlockNBT ..
type StructureBlockNBT struct {
	AnimationMode    byte    `mapstructure:"animationMode"`
	AnimationSeconds float32 `mapstructure:"animationSeconds"`
	Data             int32   `mapstructure:"data"`
	DataField        string  `mapstructure:"dataField"`
	IgnoreEntities   byte    `mapstructure:"ignoreEntities"`
	IncludePlayers   byte    `mapstructure:"includePlayers"`
	Integrity        float32 `mapstructure:"integrity"`
	Mirror           byte    `mapstructure:"mirror"`
	RedstoneSaveMode int32   `mapstructure:"redstoneSaveMode"`
	RemoveBlocks     byte    `mapstructure:"removeBlocks"`
	Rotation         byte    `mapstructure:"rotation"`
	Seed             int64   `mapstructure:"seed"`
	ShowBoundingBox  byte    `mapstructure:"showBoundingBox"`
	StructureName    string  `mapstructure:"structureName"`
	XStructureOffset int32   `mapstructure:"xStructureOffset"`
	XStructureSize   int32   `mapstructure:"xStructureSize"`
	YStructureOffset int32   `mapstructure:"yStructureOffset"`
	YStructureSize   int32   `mapstructure:"xStructureSize"`
	ZStructureOffset int32   `mapstructure:"zStructureOffset"`
	ZStructureSize   int32   `mapstructure:"xStructureSize"`
}

// 结构方块
type StructureBlock struct {
	DefaultBlock
	NBT StructureBlockNBT
}

func (s StructureBlock) NeedSpecialHandle() bool {
	if s.NBT.AnimationMode != 0 {
		return true
	}
	if s.NBT.AnimationSeconds != 0 {
		return true
	}
	if s.NBT.Data != 1 {
		return true
	}
	if len(s.NBT.DataField) > 0 {
		return true
	}
	if s.NBT.IgnoreEntities == 1 {
		return true
	}
	if s.NBT.IncludePlayers == 1 {
		return true
	}
	if s.NBT.Integrity != 100 {
		return true
	}
	if s.NBT.Mirror == 1 {
		return true
	}
	if s.NBT.RedstoneSaveMode != 0 {
		return true
	}
	if s.NBT.RemoveBlocks == 1 {
		return true
	}
	if s.NBT.Rotation != 0 {
		return true
	}
	if s.NBT.Seed != 0 {
		return true
	}
	if s.NBT.ShowBoundingBox == 0 {
		return true
	}
	if len(s.NBT.StructureName) > 0 {
		return true
	}
	if s.NBT.XStructureOffset != 0 {
		return true
	}
	if s.NBT.XStructureSize != 5 {
		return true
	}
	if s.NBT.YStructureOffset != -1 {
		return true
	}
	if s.NBT.YStructureSize != 5 {
		return true
	}
	if s.NBT.ZStructureOffset != 0 {
		return true
	}
	if s.NBT.ZStructureSize != 5 {
		return true
	}
	return false
}

func (s StructureBlock) NeedCheckCompletely() bool {
	return false
}

func (s *StructureBlock) Parse(nbtMap map[string]any) error {
	var result StructureBlockNBT
	err := mapstructure.Decode(&nbtMap, &result)
	if err != nil {
		return fmt.Errorf("Parse: %v", err)
	}
	s.NBT = result
	return nil
}

func (s StructureBlock) StableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)
	basicInfo := s.DefaultBlock.StableBytes()

	w.ByteSlice(&basicInfo)
	w.Uint8(&s.NBT.AnimationMode)
	w.Float32(&s.NBT.AnimationSeconds)
	w.Varint32(&s.NBT.Data)
	w.String(&s.NBT.DataField)
	w.Uint8(&s.NBT.IgnoreEntities)
	w.Uint8(&s.NBT.IncludePlayers)
	w.Float32(&s.NBT.Integrity)
	w.Uint8(&s.NBT.Mirror)
	w.Varint32(&s.NBT.RedstoneSaveMode)
	w.Uint8(&s.NBT.RemoveBlocks)
	w.Uint8(&s.NBT.Rotation)
	w.Int64(&s.NBT.Seed)
	w.Uint8(&s.NBT.ShowBoundingBox)
	w.String(&s.NBT.StructureName)
	w.Varint32(&s.NBT.XStructureOffset)
	w.Varint32(&s.NBT.XStructureSize)
	w.Varint32(&s.NBT.YStructureOffset)
	w.Varint32(&s.NBT.YStructureSize)
	w.Varint32(&s.NBT.ZStructureOffset)
	w.Varint32(&s.NBT.ZStructureSize)

	return buf.Bytes()
}
