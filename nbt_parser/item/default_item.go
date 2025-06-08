package nbt_parser

import (
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
)

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
	if block.SubBlock != nil {
		if block.SubBlock.NeedSpecialHandle() {
			shouldCleanItemLock = true
		}
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
	block, err := ParseItemBlock(nbtMap)
	if err != nil {
		return fmt.Errorf("ParseNormal: %v", err)
	}
	// Parse data
	d.parse(basic, enhance, block)
	// Return
	return nil
}

func (d *DefaultItem) ParseNetwork(item protocol.ItemStack, itemNetworkIDToName map[int32]string) error {
	// Parse basic item data
	basic, err := ParseItemBasicDataNetwork(item, itemNetworkIDToName)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	// Parse enhance item data
	enhance, err := ParseItemEnhanceNetwork(item)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	// Parse item block data
	block, err := ParseItemBlockNetwork(item)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	// Parse data
	d.parse(basic, enhance, block)
	// Return
	return nil
}

func (DefaultItem) NeedSpecialHandle() bool {
	panic("TODO")
}

func (DefaultItem) NeedCheckCompletely() bool {
	panic("TODO")
}

func (DefaultItem) TypeStableBytes() []byte {
	panic("TODO")
}

func (DefaultItem) FullStableBytes() []byte {
	panic("TODO")
}
