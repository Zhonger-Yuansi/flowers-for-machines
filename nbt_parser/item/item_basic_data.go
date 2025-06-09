package nbt_parser_item

import (
	"fmt"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
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

	tag, ok := nbtMap["tag"].(map[string]any)
	if ok {
		damage, ok := tag["Damage"].(int32)
		if ok {
			result.Metadata = int16(damage)
		}
	}

	return result, nil
}

func ParseItemBasicDataNetwork(item protocol.ItemStack, itemName string) (result ItemBasicData, err error) {
	result.Name = strings.ToLower(itemName)
	if !strings.HasPrefix(result.Name, "minecraft:") {
		result.Name = "minecraft:" + result.Name
	}

	result.Count = uint8(item.Count)
	result.Metadata = int16(item.MetadataValue)

	if item.NBTData != nil {
		damage, ok := item.NBTData["Damage"].(int32)
		if ok {
			result.Metadata = int16(damage)
		}
	}

	return
}
