package item_stack_transaction

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface/item_stack_operation"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// ItemStackTransaction 是单个物品操作事务，
// 它希望使用者尽可能多的将物品堆栈请求内联在一个数据包中，
// 这样可以有效的节省操作的时间消耗
type ItemStackTransaction struct {
	api            *resources_control.Resources
	operations     []item_stack_operation.ItemStackOperation
	stackNetworkID map[resources_control.SlotLocation]int32
}

// NewItemStackTransaction 基于 api 创建并返回一个新的 ItemStackTransaction
func NewItemStackTransaction(api *resources_control.Resources) *ItemStackTransaction {
	return &ItemStackTransaction{
		api:            api,
		operations:     nil,
		stackNetworkID: make(map[resources_control.SlotLocation]int32),
	}
}
