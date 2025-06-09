package nbt_parser_item

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/TriM-Organization/bedrock-world-operator/block"
	"github.com/df-mc/worldupgrader/blockupgrader"
)

type ItemEnhanceData struct {
	ItemComponent utils.ItemComponent
	DisplayName   string
	EnchList      []SingleItemEnch
}

// Marshal ..
func (i *ItemEnhanceData) Marshal(io protocol.IO) {
	protocol.Single(io, &i.ItemComponent)
	io.String(&i.DisplayName)
	protocol.SliceUint16Length(io, &i.EnchList)
}

func ParseItemEnhance(nbtMap map[string]any) (result ItemEnhanceData, err error) {
	result.ItemComponent = utils.ParseItemComponent(nbtMap)

	result.EnchList, err = ParseItemEnchList(nbtMap)
	if err != nil {
		return result, fmt.Errorf("ParseItemEnhance: %v", err)
	}

	tag, ok := nbtMap["tag"].(map[string]any)
	if !ok {
		return
	}
	display, ok := tag["display"].(map[string]any)
	if !ok {
		return
	}
	result.DisplayName, _ = display["Name"].(string)

	return
}

func ParseItemEnhanceNetwork(item protocol.ItemStack) (result ItemEnhanceData, err error) {
	result.ItemComponent = utils.ParseItemComponentNetwork(item)

	result.EnchList, err = ParseItemEnchListNetwork(item)
	if err != nil {
		return result, fmt.Errorf("ParseItemEnhanceNetwork: %v", err)
	}

	if item.NBTData == nil {
		return
	}
	display, ok := item.NBTData["display"].(map[string]any)
	if !ok {
		return
	}
	result.DisplayName, _ = display["Name"].(string)

	return
}

type ItemBlockData struct {
	Name     string
	States   map[string]any
	SubBlock nbt_parser_interface.Block
}

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
