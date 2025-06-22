package resources_control

import (
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

// ConstantPacket 记载在登录序列期间，
// 由租赁服发送的在整个连接期间不会变化的常量
type ConstantPacket struct {
	// 所有可用物品
	availableItems       []protocol.ItemEntry
	itemNetworkIDMapping map[int32]int
	itemNameMapping      map[string]int
	itemNameMappingInv   []string
	// 创造物品
	creativeContent    []protocol.CreativeItem
	creativeNIMapping  map[int32][]int // NI: Network ID
	creativeCNIMapping map[uint32]int  // CNI: Creative Network ID
	// 所有可通过指令获得的物品
	commandItems        []string
	commandItemsMapping map[string]bool
}

// NewConstantPacket 创建并返回一个新的 ConstantPacket
func NewConstantPacket() *ConstantPacket {
	return &ConstantPacket{
		availableItems:       nil,
		itemNetworkIDMapping: make(map[int32]int),
		itemNameMapping:      make(map[string]int),
		itemNameMappingInv:   nil,
		creativeContent:      nil,
		creativeNIMapping:    make(map[int32][]int),
		creativeCNIMapping:   make(map[uint32]int),
		commandItems:         nil,
		commandItemsMapping:  make(map[string]bool),
	}
}

// AllCreativeContent 返回租赁服在登录序列发送的创造物品数据。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) AllCreativeContent() []protocol.CreativeItem {
	return c.creativeContent
}

// CreativeItemByCNI 返回创造物品网络 ID 为 creativeNetworkID 的创造物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) CreativeItemByCNI(creativeNetworkID uint32) protocol.CreativeItem {
	return c.creativeContent[c.creativeCNIMapping[creativeNetworkID]]
}

// CreativeItemByNI 返回物品数字网络 ID 为 networkID 的多个创造物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) CreativeItemByNI(networkID int32) []protocol.CreativeItem {
	result := make([]protocol.CreativeItem, 0)
	for _, index := range c.creativeNIMapping[networkID] {
		result = append(result, c.creativeContent[index])
	}
	return result
}

// CreativeItemByName 返回名称为 name 的多个创造物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) CreativeItemByName(name string) []protocol.CreativeItem {
	name = strings.ToLower(name)
	if !strings.HasPrefix(name, "minecraft:") {
		name = "minecraft:" + name
	}
	return c.CreativeItemByNI(int32(c.ItemByName(name).RuntimeID))
}

// onCreativeContent ..
func (c *ConstantPacket) onCreativeContent(p *packet.CreativeContent) {
	c.creativeContent = p.Items
	for index, item := range p.Items {
		c.creativeNIMapping[item.Item.NetworkID] = append(c.creativeNIMapping[item.Item.NetworkID], index)
		c.creativeCNIMapping[item.CreativeItemNetworkID] = index
	}
}

// AllAvailableItems 返回租赁服在登录序列发送的所有可用物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) AllAvailableItems() []protocol.ItemEntry {
	return c.availableItems
}

// ItemByNetworkID 返回网络 ID 为 networkID 的物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) ItemByNetworkID(networkID int32) protocol.ItemEntry {
	return c.availableItems[c.itemNetworkIDMapping[networkID]]
}

// ItemByName 返回物品名称为 name 的物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) ItemByName(name string) protocol.ItemEntry {
	name = strings.ToLower(name)
	if !strings.HasPrefix(name, "minecraft:") {
		name = "minecraft:" + name
	}
	return c.availableItems[c.itemNameMapping[name]]
}

// ItemNameByNetworkID 返回网络 ID 为 networkID 的物品的名称
func (c ConstantPacket) ItemNameByNetworkID(networkID int32) string {
	return c.itemNameMappingInv[c.itemNetworkIDMapping[networkID]]
}

// updateByGameData ..
func (c *ConstantPacket) updateByGameData(data minecraft.GameData) {
	c.availableItems = data.Items
	c.itemNameMappingInv = make([]string, len(c.availableItems))
	for index, item := range c.availableItems {
		c.itemNetworkIDMapping[int32(item.RuntimeID)] = index
		c.itemNameMapping[item.Name] = index
		c.itemNameMappingInv[index] = item.Name
	}
}

// AllCommandItems 返回可以通过指令获得的全部物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) AllCommandItems() []string {
	return c.commandItems
}

// ItemCanGetByCommand 检查物品名为 name 的物品是否可以通过命令获取
func (c ConstantPacket) ItemCanGetByCommand(name string) bool {
	name = strings.ToLower(name)
	if !strings.HasPrefix(name, "minecraft:") {
		name = "minecraft:" + name
	}

	if name == "minecraft:written_book" {
		return true
	}

	return c.commandItemsMapping[name]
}

// onAvailableCommands ..
func (c *ConstantPacket) onAvailableCommands(p *packet.AvailableCommands) {
	c.commandItems = nil
	c.commandItemsMapping = make(map[string]bool)

	for _, enum := range p.Enums {
		if enum.Type != "Item" {
			continue
		}

		for _, index := range enum.ValueIndices {
			itemName := p.EnumValues[index]
			if !strings.HasPrefix(itemName, "minecraft:") {
				continue
			}
			c.commandItems = append(c.commandItems, itemName)
			c.commandItemsMapping[itemName] = true
		}

		return
	}

	panic("onAvailableCommands: Should nerver happened")
}
