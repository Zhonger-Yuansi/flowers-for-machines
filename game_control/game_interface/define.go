package game_interface

import "github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"

// 描述机器人的基本信息
type BotBasicInfo struct {
	BotName         string // 机器人名称
	XUID            string // 机器人 XUID
	EntityUniqueID  int64  // 机器人唯一 ID
	EntityRuntimeID uint64 // 机器人运行时 ID
}

// 用于 PhoenixBuilder 与租赁服交互。
// 此结构体下的实现将允许您与租赁服进行交互操作，例如打开容器等
type GameInterface struct {
	// 用于向租赁服发送数据包的函数
	WritePacket func(packet.Packet) error
	// 存储客户端的基本信息
	BotBasicInfo BotBasicInfo
}
