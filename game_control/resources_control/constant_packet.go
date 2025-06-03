package resources_control

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

// ConstantPacket 记载在登录序列期间，
// 由租赁服发送的在整个连接期间不会变化的常量
type ConstantPacket struct {
	creativeContent    []protocol.CreativeItem
	creativeNIMapping  map[int32]int  // NI: Network ID
	creativeCNIMapping map[uint32]int // CNI: Creative Network ID
}

// NewConstantPacket 创建并返回一个新的 ConstantPacket
func NewConstantPacket() *ConstantPacket {
	return &ConstantPacket{
		creativeContent:    nil,
		creativeNIMapping:  make(map[int32]int),
		creativeCNIMapping: make(map[uint32]int),
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

// CreativeItemByNI 返回物品数字网络 ID 为 networkID 的创造物品。
// 使用者不应修改返回的值，否则不保证程序的行为是正确的
func (c ConstantPacket) CreativeItemByNI(networkID int32) protocol.CreativeItem {
	return c.creativeContent[c.creativeNIMapping[networkID]]
}

// onCreativeContent ..
func (c *ConstantPacket) onCreativeContent(p *packet.CreativeContent) {
	c.creativeContent = p.Items
	for index, item := range p.Items {
		c.creativeNIMapping[item.Item.NetworkID] = index
		c.creativeCNIMapping[item.CreativeItemNetworkID] = index
	}
}
