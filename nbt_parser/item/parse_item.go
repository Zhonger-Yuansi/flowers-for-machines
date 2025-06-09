package nbt_parser

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
)

func ParseNBTItemNormal(nbtMap map[string]any) (item nbt_parser_interface.Item, err error) {
	var defaultItem DefaultItem

	err = defaultItem.ParseNormal(nbtMap)
	if err != nil {
		return nil, fmt.Errorf("ParseNBTItemNormal: %v", err)
	}

	itemType, ok := mapping.SupportItemsPool[defaultItem.ItemName()]
	if !ok {
		return &defaultItem, nil
	}

	switch itemType {
	case mapping.SupportNBTItemTypeBook:
		item = &Book{DefaultItem: defaultItem}
	case mapping.SupportNBTItemTypeBanner:
		item = &Banner{DefaultItem: defaultItem}
	case mapping.SupportNBTItemTypeShield:
		item = &Shield{DefaultItem: defaultItem}
	default:
		panic("ParseNBTItemNormal: Should nerver happened")
	}

	err = item.ParseNormal(nbtMap)
	if err != nil {
		return nil, fmt.Errorf("ParseNBTItemNormal: %v", err)
	}
	return item, nil
}

func ParseNBTItemNetwork(itemStack protocol.ItemStack, itemNetworkIDToName map[int32]string) (item nbt_parser_interface.Item, err error) {
	var defaultItem DefaultItem

	err = defaultItem.ParseNetwork(itemStack, itemNetworkIDToName)
	if err != nil {
		return nil, fmt.Errorf("ParseNBTItemNetwork: %v", err)
	}

	itemType, ok := mapping.SupportItemsPool[defaultItem.ItemName()]
	if !ok {
		return &defaultItem, nil
	}

	switch itemType {
	case mapping.SupportNBTItemTypeBook:
		item = &Book{DefaultItem: defaultItem}
	case mapping.SupportNBTItemTypeBanner:
		item = &Banner{DefaultItem: defaultItem}
	case mapping.SupportNBTItemTypeShield:
		item = &Shield{DefaultItem: defaultItem}
	default:
		panic("ParseNBTItemNetwork: Should nerver happened")
	}

	err = item.ParseNetwork(itemStack, itemNetworkIDToName)
	if err != nil {
		return nil, fmt.Errorf("ParseNBTItemNetwork: %v", err)
	}
	return item, nil
}

func init() {
	nbt_parser_interface.ParseNBTItemNormal = ParseNBTItemNormal
	nbt_parser_interface.ParseNBTItemNetwork = ParseNBTItemNetwork
}
