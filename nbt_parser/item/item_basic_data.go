package nbt_parser

import (
	"fmt"

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
