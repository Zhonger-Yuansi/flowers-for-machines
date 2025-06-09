package nbt_parser_item

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/TriM-Organization/bedrock-world-operator/block"
	"github.com/df-mc/worldupgrader/blockupgrader"
)

// ItemBlockData 指示该物品是一个方块，
// 或该物品可以作为一个方块进行放置
type ItemBlockData struct {
	// Name 是这个方块的名称
	Name string
	// States 是这个方块的方块状态
	States map[string]any
	// SubBlock 是这个方块的附加数据。
	// 如果这个方块是已被支持的 NBT 方块，
	// 且 NBT 字段存在，则 SubBlock 非空
	SubBlock nbt_parser_interface.Block
}

// ParseItemBlock ..
func ParseItemBlock(itemName string, nbtMap map[string]any) (result ItemBlockData, err error) {
	var blockMap map[string]any
	var haveBlock bool

	blockMap, haveBlock = nbtMap["Block"].(map[string]any)
	tag, _ := nbtMap["tag"].(map[string]any)

	if haveBlock {
		name, _ := blockMap["name"].(string)
		states, _ := blockMap["states"].(map[string]any)
		version, _ := blockMap["version"].(int32)

		newBlock := blockupgrader.Upgrade(blockupgrader.BlockState{
			Name:       name,
			Properties: states,
			Version:    version,
		})
		result.Name = newBlock.Name
		result.States = newBlock.Properties
	} else {
		blockName, ok := mapping.ItemNameToBlockName[itemName]
		if !ok {
			return
		}
		blockType, ok := mapping.SupportBlocksPool[blockName]
		if !ok {
			panic("ParseItemBlock: Should nerver happened")
		}
		if !mapping.SubBlocksPool[blockType] {
			return
		}

		rid, found := block.StateToRuntimeID(blockName, map[string]any{})
		if !found {
			panic("ParseItemBlock: Should nerver happened")
		}

		name, states, found := block.RuntimeIDToState(rid)
		if !found {
			panic("ParseItemBlock: Should nerver happened")
		}

		result.Name = name
		result.States = states
	}

	if len(tag) > 0 {
		result.SubBlock, err = nbt_parser_interface.ParseBlock(result.Name, result.States, tag)
		if err != nil {
			return ItemBlockData{}, fmt.Errorf("ParseItemBlock: %v", err)
		}
	}

	return
}

// ParseItemBlockNetwork ..
func ParseItemBlockNetwork(itemName string, item protocol.ItemStack) (result ItemBlockData, err error) {
	if item.BlockRuntimeID == 0 {
		name, states, found := block.RuntimeIDToState(uint32(item.BlockRuntimeID))
		if !found {
			panic("ParseItemBlockNetwork: Should nerver happened")
		}
		result.Name = name
		result.States = states
	} else {
		blockName, ok := mapping.ItemNameToBlockName[itemName]
		if !ok {
			return
		}
		blockType, ok := mapping.SupportBlocksPool[blockName]
		if !ok {
			panic("ParseItemBlockNetwork: Should nerver happened")
		}
		if !mapping.SubBlocksPool[blockType] {
			return
		}

		rid, found := block.StateToRuntimeID(blockName, map[string]any{})
		if !found {
			panic("ParseItemBlockNetwork: Should nerver happened")
		}

		name, states, found := block.RuntimeIDToState(rid)
		if !found {
			panic("ParseItemBlockNetwork: Should nerver happened")
		}

		result.Name = name
		result.States = states
	}

	if len(item.NBTData) > 0 {
		result.SubBlock, err = nbt_parser_interface.ParseBlock(result.Name, result.States, item.NBTData)
		if err != nil {
			return ItemBlockData{}, fmt.Errorf("ParseItemBlockNetwork: %v", err)
		}
	}

	return
}
