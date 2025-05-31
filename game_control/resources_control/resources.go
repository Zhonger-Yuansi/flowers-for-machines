package resources_control

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/client"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/google/uuid"
)

type Resources struct {
	// client 是连接到租赁服的基本客户端
	client *client.Client

	// commandCallback 存放所有命令请求的回调函数
	commandCallback utils.SyncMap[uuid.UUID, func(pk *packet.CommandOutput)]

	// inventory 持有机器人已经拥有或打开的库存
	inventory Inventories
	// inventoryCallback 存放所有物品更改的回调函数
	inventoryCallback utils.SyncMap[SlotLocation, utils.MultipleCallback]

	itemStackMapping  utils.SyncMap[ItemStackRequestID, ItemStackResponseMapping]
	itemStackUpdater  utils.SyncMap[ItemStackRequestID, map[SlotLocation]ExpectedNewItem]
	itemStackCallback utils.SyncMap[ItemStackRequestID, func()]
	// // 管理物品操作请求及结果
	// ItemStackOperation item_stack_request_with_response
	// // 管理容器资源的占用状态，同时存储容器操作的结果
	// Container container
	// // 管理结构资源并保存结构请求的回应
	// Structure mcstructure
	// // 数据包监听器
	// Listener packet_listener
	// // 管理和保存其他小型的资源，
	// // 例如游戏刻相关
	// Others others
}
