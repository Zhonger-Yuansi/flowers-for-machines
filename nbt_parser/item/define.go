package nbt_parser

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/TriM-Organization/bedrock-world-operator/block"
	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/df-mc/worldupgrader/itemupgrader"
	"github.com/mitchellh/mapstructure"
)

type ItemBasicData struct {
	Name     string `mapstructure:"Name"`
	Count    uint8  `mapstructure:"Count"`
	Metadata int16  `mapstructure:"Damage"`
}

func ParseItemBasicData(nbtMap map[string]any) (result ItemBasicData, err error) {
	err = mapstructure.Decode(&nbtMap, &result)
	if err != nil {
		return result, fmt.Errorf("ParseItemBasicData: %v", err)
	}

	newItem := itemupgrader.Upgrade(itemupgrader.ItemMeta{
		Name: result.Name,
		Meta: result.Metadata,
	})
	result.Name = newItem.Name
	result.Metadata = newItem.Meta

	return result, nil
}

func ParseItemBasicDataNetwork(item protocol.ItemStack, itemNetworkIDToName map[int32]string) (result ItemBasicData, err error) {
	result.Name = itemNetworkIDToName[item.ItemType.NetworkID]
	result.Count = uint8(item.Count)
	result.Metadata = int16(item.MetadataValue)

	if len(result.Name) == 0 {
		return ItemBasicData{}, fmt.Errorf(
			"ParseItemBasicDataNetwork: itemNetworkIDToName not record the name of item network ID %d; itemNetworkIDToName = %#v",
			item.ItemType.NetworkID, itemNetworkIDToName,
		)
	}

	return
}

type SingleItemEnch struct {
	ID    int16 `mapstructure:"id"`
	Level int16 `mapstructure:"lvl"`
}

func parseItemEnchList(enchList []any) (result []SingleItemEnch, err error) {
	for _, value := range enchList {
		var singleItemEnch SingleItemEnch

		val, ok := value.(map[string]any)
		if !ok {
			continue
		}

		err = mapstructure.Decode(&val, &singleItemEnch)
		if err != nil {
			return nil, fmt.Errorf("ParseItemEnchList: %v", err)
		}

		result = append(result, singleItemEnch)
	}
	return
}

func ParseItemEnchList(nbtMap map[string]any) (result []SingleItemEnch, err error) {
	tag, ok := nbtMap["tag"].(map[string]any)
	if !ok {
		return
	}

	ench, ok := tag["ench"].([]any)
	if !ok {
		return
	}

	result, err = parseItemEnchList(ench)
	if err != nil {
		return nil, fmt.Errorf("ParseItemEnchList: %v", err)
	}

	return
}

func ParseItemEnchListNetwork(item protocol.ItemStack) (result []SingleItemEnch, err error) {
	if item.NBTData == nil {
		return
	}

	ench, ok := item.NBTData["ench"].([]any)
	if !ok {
		return
	}

	result, err = parseItemEnchList(ench)
	if err != nil {
		return nil, fmt.Errorf("ParseItemEnchListNetwork: %v", err)
	}

	return
}

type ItemEnhanceData struct {
	ItemComponent utils.ItemComponent
	DisplayName   string
	EnchList      []SingleItemEnch
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

func ParseItemBlock(nbtMap map[string]any) (result ItemBlockData, err error) {
	block, ok := nbtMap["Block"].(map[string]any)
	if !ok {
		return
	}

	name, _ := block["name"].(string)
	states, _ := block["states"].(map[string]any)
	version, _ := block["version"].(int32)
	tag, _ := nbtMap["tag"].(map[string]any)

	newBlock := blockupgrader.Upgrade(blockupgrader.BlockState{
		Name:       name,
		Properties: states,
		Version:    version,
	})
	result.Name = newBlock.Name
	result.States = newBlock.Properties

	if len(tag) > 0 {
		result.SubBlock, err = nbt_parser_interface.ParseBlock(result.Name, result.States, tag)
		if err != nil {
			return ItemBlockData{}, fmt.Errorf("ParseItemBlock: %v", err)
		}
	}

	return
}

func ParseItemBlockNetwork(item protocol.ItemStack) (result ItemBlockData, err error) {
	if item.BlockRuntimeID == 0 {
		return
	}

	name, states, found := block.RuntimeIDToState(uint32(item.BlockRuntimeID))
	if !found {
		panic("ParseItemBlockNetwork: Should nerver happened")
	}

	result.Name = name
	result.States = states
	tag := item.NBTData

	if len(tag) > 0 {
		result.SubBlock, err = nbt_parser_interface.ParseBlock(result.Name, result.States, tag)
		if err != nil {
			return ItemBlockData{}, fmt.Errorf("ParseItemBlockNetwork: %v", err)
		}
	}

	return
}
