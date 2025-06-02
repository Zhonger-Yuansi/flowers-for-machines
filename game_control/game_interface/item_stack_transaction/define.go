package item_stack_transaction

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
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

// Discord 丢弃曾经执行的更改。
// 从本质上说，它将清空底层操作序列
func (i *ItemStackTransaction) Discord() {
	i.operations = nil
}

// Commit 将底层操作序列内联到单个的物品堆栈请求的数据包中，并将它发送至租赁服。
// 如果没有返回错误，Commit 在完成后将使用 Discord 清空底层操作序列
//
// Commit 在设计上考虑并预期事务的所有都会成功，因此内联将尽可能紧凑，而这依赖于“成功”的预期前提。
// 这意味着，一旦某个步骤失败，那么整个物品堆栈操作都可能失败，并且最终的结果将是未定义的。
//
// success 为真指示该事务的全部操作完全成功，若为否则可能部分失败；
// pk 指示最终编译得到的数据包，它可以用于调试，但不应重新用于发送；
// serverResponse 则指示租赁服针对 pk 的响应结果
func (i *ItemStackTransaction) Commit() (
	success bool,
	pk *packet.ItemStackRequest,
	serverResponse []*protocol.ItemStackResponse,
	err error,
) {
	api := i.api

	pk = new(packet.ItemStackRequest)
	handler := newItemStackOperationHandler(
		api.Container(),
		newVirtualInventories(api.Inventories()),
		newResponseMapping(),
	)
	requests := make([][]item_stack_operation.ItemStackOperation, 0)
	waiters := make([]chan struct{}, 0)
	requestIDs := make([]resources_control.ItemStackRequestID, 0)

	// Step 1: Split by operations that can't inline
	currentRequest := make([]item_stack_operation.ItemStackOperation, 0)
	for _, operation := range i.operations {
		if !operation.CanInline() {
			if len(currentRequest) != 0 {
				requests = append(requests, currentRequest)
			}
			requests = append(requests, []item_stack_operation.ItemStackOperation{operation})
			currentRequest = nil
			continue
		}
		currentRequest = append(currentRequest, operation)
	}
	if len(currentRequest) != 0 {
		requests = append(requests, currentRequest)
		currentRequest = nil
	}

	// Step 2: Get new request ID for all operations
	serverResponse = make([]*protocol.ItemStackResponse, len(requests))
	for range requests {
		requestIDs = append(requestIDs, api.ItemStackOperation().NewRequestID())
	}

	// Step 3: Construct stack request actions
	for index, request := range requests {
		var itemNewName *string
		var updater map[resources_control.SlotLocation]resources_control.ExpectedNewItem

		actions := make([]protocol.StackRequestAction, 0)
		requestID := requestIDs[index]

		for _, operation := range request {
			var result []protocol.StackRequestAction
			var err error

			switch op := operation.(type) {
			case item_stack_operation.Move:
				result, err = handler.handleMove(op, requestID)
			case item_stack_operation.Swap:
				result, err = handler.handleSwap(op, requestID)
			case item_stack_operation.Drop:
				result, err = handler.handleDrop(op, requestID)
			case item_stack_operation.DropHotbar:
				result, err = handler.handleDropHotbar(op, requestID)
			case item_stack_operation.CreativeItem:
				result, err = handler.handleCreativeItem(op, requestID)
			case item_stack_operation.Renaming:
				result, err = handler.handleRenaming(op, requestID)
				itemNewName = &op.NewName
			case item_stack_operation.Looming:
				result, err = handler.handleLooming(op, requestID)
				updater = make(map[resources_control.SlotLocation]resources_control.ExpectedNewItem)
				updater[op.BannerPath] = op.ResultItem
			}

			if err != nil {
				return false, nil, nil, fmt.Errorf("Commit: %v", err)
			}

			actions = append(actions, result...)
		}

		newRequest := protocol.ItemStackRequest{
			RequestID: int32(requestID),
			Actions:   actions,
		}
		if itemNewName != nil {
			newRequest.FilterStrings = []string{*itemNewName}
			newRequest.FilterCause = protocol.FilterCauseAnvilText
		}
		pk.Requests = append(pk.Requests, newRequest)

		idx := index
		channel := make(chan struct{})
		waiters = append(waiters, channel)

		api.ItemStackOperation().AddNewRequest(
			requestID,
			handler.responseMapping.mapping,
			updater,
			func(response *protocol.ItemStackResponse) {
				serverResponse[idx] = response
				close(channel)
			},
		)
	}

	// Step 4: Send packet
	err = api.WritePacket(pk)
	if err != nil {
		return false, nil, nil, fmt.Errorf("Commit: %v", err)
	}

	// Step 5: Waiting all finished
	for _, waiter := range waiters {
		<-waiter
	}

	// Setp 6: Return unsuccess
	for _, response := range serverResponse {
		if response.Status != protocol.ItemStackResponseStatusOK {
			i.Discord()
			return false, pk, serverResponse, nil
		}
	}

	// Step 7: Return success
	i.Discord()
	return true, pk, serverResponse, nil
}
