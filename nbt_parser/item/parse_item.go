package nbt_parser_item

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
)

// ParseItemNormal 从 nbtMap 解析一个 NBT 物品。
// nbtMap 是含有这个物品 tag 标签的父复合标签
func ParseItemNormal(nbtMap map[string]any) (item nbt_parser_interface.Item, err error) {
	var defaultItem DefaultItem

	err = defaultItem.ParseNormal(nbtMap)
	if err != nil {
		return nil, fmt.Errorf("ParseItemNormal: %v", err)
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
		panic("ParseItemNormal: Should nerver happened")
	}

	err = item.ParseNormal(nbtMap)
	if err != nil {
		return nil, fmt.Errorf("ParseItemNormal: %v", err)
	}
	return item, nil
}

// ParseItemNetwork 解析网络传输上的物品堆栈实例 item。
// itemName 是这个物品堆栈实例的名称
func ParseItemNetwork(itemStack protocol.ItemStack, itemName string) (item nbt_parser_interface.Item, err error) {
	var defaultItem DefaultItem

	err = defaultItem.ParseNetwork(itemStack, itemName)
	if err != nil {
		return nil, fmt.Errorf("ParseItemNetwork: %v", err)
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
		panic("ParseItemNetwork: Should nerver happened")
	}

	err = item.ParseNetwork(itemStack, itemName)
	if err != nil {
		return nil, fmt.Errorf("ParseItemNetwork: %v", err)
	}
	return item, nil
}

func init() {
	nbt_parser_interface.ParseItemNormal = ParseItemNormal
	nbt_parser_interface.ParseItemNetwork = ParseItemNetwork
}
