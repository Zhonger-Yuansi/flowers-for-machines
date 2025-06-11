package nbt_parser_item

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// 默认 NBT 物品
type DefaultItem struct {
	Basic   ItemBasicData
	Enhance ItemEnhanceData
	Block   ItemBlockData
}

func (d *DefaultItem) ItemName() string {
	d.Basic.Name = strings.ToLower(d.Basic.Name)
	if !strings.HasPrefix(d.Basic.Name, "minecraft:") {
		d.Basic.Name = "minecraft:" + d.Basic.Name
	}
	return d.Basic.Name
}

func (d DefaultItem) ItemCount() uint8 {
	return d.Basic.Count
}

func (d DefaultItem) ItemMetadata() int16 {
	return d.Basic.Metadata
}

func (d *DefaultItem) parse(basic ItemBasicData, enhance ItemEnhanceData, block ItemBlockData) {
	// Prepare
	var shouldCleanItemLock bool
	// Fix logic problem
	if len(block.Name) != 0 {
		enhance.EnchList = nil
	}
	if block.SubBlock != nil && block.SubBlock.NeedSpecialHandle() {
		shouldCleanItemLock = true
	}
	if len(enhance.EnchList) > 0 || len(enhance.DisplayName) > 0 {
		shouldCleanItemLock = true
	}
	if shouldCleanItemLock {
		enhance.ItemComponent.LockInInventory = false
		enhance.ItemComponent.LockInSlot = false
	}
	// Sync data
	*d = DefaultItem{
		Basic:   basic,
		Enhance: enhance,
		Block:   block,
	}
}

func (d *DefaultItem) ParseNormal(nbtMap map[string]any) error {
	// Parse basic item data
	basic, err := ParseItemBasicData(nbtMap)
	if err != nil {
		return fmt.Errorf("ParseNormal: %v", err)
	}
	// Parse enhance item data
	enhance, err := ParseItemEnhance(nbtMap)
	if err != nil {
		return fmt.Errorf("ParseNormal: %v", err)
	}
	// Parse item block data
	block, err := ParseItemBlock(basic.Name, nbtMap)
	if err != nil {
		return fmt.Errorf("ParseNormal: %v", err)
	}
	// Parse data
	d.parse(basic, enhance, block)
	// Return
	return nil
}

func (d *DefaultItem) ParseNetwork(item protocol.ItemStack, itemName string) error {
	// Parse basic item data
	basic, err := ParseItemBasicDataNetwork(item, itemName)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	// Parse enhance item data
	enhance, err := ParseItemEnhanceNetwork(item)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	// Parse item block data
	block, err := ParseItemBlockNetwork(basic.Name, item)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	// Parse data
	d.parse(basic, enhance, block)
	// Return
	return nil
}

func (d *DefaultItem) UnderlyingItem() nbt_parser_interface.Item {
	return d
}

func (d DefaultItem) NeedEnchOrRename() bool {
	if len(d.Enhance.DisplayName) > 0 || len(d.Enhance.EnchList) > 0 {
		return true
	}
	if d.Block.SubBlock != nil && d.Block.SubBlock.NeedSpecialHandle() {
		return true
	}
	return false
}

func (d DefaultItem) IsComplex() bool {
	if d.Block.SubBlock != nil && d.Block.SubBlock.NeedSpecialHandle() {
		return true
	}
	return false
}

func (d DefaultItem) NeedCheckCompletely() bool {
	if len(d.Enhance.EnchList) > 0 {
		return true
	}
	if d.Block.SubBlock != nil && d.Block.SubBlock.NeedCheckCompletely() {
		return true
	}
	return false
}

func (d DefaultItem) NBTStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)
	haveBlock := (len(d.Block.Name) > 0)

	// ItemComponent
	protocol.Single(w, &d.Enhance.ItemComponent)

	// Block
	w.Bool(&haveBlock)
	if haveBlock {
		blockStatesString := utils.MarshalBlockStates(d.Block.States)
		haveSubBlock := (d.Block.SubBlock != nil)

		w.String(&d.Block.Name)
		w.String(&blockStatesString)

		w.Bool(&haveSubBlock)
		if haveSubBlock {
			subBlockData := d.Block.SubBlock.StableBytes()
			w.ByteSlice(&subBlockData)
		}
	}

	return buf.Bytes()
}

func (d *DefaultItem) TypeStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)
	name := d.ItemName()

	// Basic
	w.String(&name)
	w.Int16(&d.Basic.Metadata)

	// Enhance (Display Name, Ench List)
	w.String(&d.Enhance.DisplayName)
	protocol.SliceUint16Length(w, &d.Enhance.EnchList)

	return append(d.NBTStableBytes(), buf.Bytes()...)
}

func (d *DefaultItem) FullStableBytes() []byte {
	return append(d.TypeStableBytes(), d.Basic.Count)
}
