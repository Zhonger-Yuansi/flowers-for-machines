package item_stack_transaction

import (
	"fmt"
	"sync"

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

// Commit 将底层操作序列内联并使用尽可能少的物品堆栈请求的数据包执行物品堆栈操作事务。
// 如果没有返回错误，Commit 在完成后将使用 Discord 清空底层操作序列。
//
// Commit 在设计上考虑并预期事务的所有都会成功，因此内联将尽可能紧凑，而这依赖于“成功”的预期前提。
// 这意味着，一旦某个步骤失败，那么整个物品堆栈操作都可能失败，并且最终的结果将是未定义的。
//
// success 为真指示该事务的全部操作完全成功，若为否则可能部分失败。
// 作为一种特殊情况，如果底层操作序列为空，则 success 总是真。
//
// pks 指示最终编译得到的多个数据包，它可以用于调试，但不应重新用于发送；
// serverResponse 则指示租赁服针对 pks 中每个物品堆栈请求的结果
func (i *ItemStackTransaction) Commit() (
	success bool,
	pks []*packet.ItemStackRequest,
	serverResponse [][]*protocol.ItemStackResponse,
	err error,
) {
	if len(i.operations) == 0 {
		return true, nil, make([][]*protocol.ItemStackResponse, 0), nil
	}

	api := i.api
	mu := new(sync.Mutex)
	requests := make([][]item_stack_operation.ItemStackOperation, 0)

	// Step 1: Split inline and can't inline
	// e.g.
	//		[
	// 			[inline, inline, ...],
	// 			[can't inline, can't inline, ...],
	// 			...
	//		]
	{
		canInlineRequests := make([]item_stack_operation.ItemStackOperation, 0)
		canNotInlineRequests := make([]item_stack_operation.ItemStackOperation, 0)

		for _, operation := range i.operations {
			if !operation.CanInline() {
				if len(canInlineRequests) > 0 {
					requests = append(requests, canInlineRequests)
					canInlineRequests = nil
				}
				canNotInlineRequests = append(canNotInlineRequests, operation)
				continue
			}
			if len(canNotInlineRequests) > 0 {
				requests = append(requests, canNotInlineRequests)
				canNotInlineRequests = nil
			}
			canInlineRequests = append(canInlineRequests, operation)
		}

		if len(canInlineRequests) != 0 {
			requests = append(requests, canInlineRequests)
			canInlineRequests = nil
		}
		if len(canNotInlineRequests) != 0 {
			requests = append(requests, canNotInlineRequests)
			canNotInlineRequests = nil
		}

		serverResponse = make([][]*protocol.ItemStackResponse, len(requests))
	}

	// Step 3: Commit
	for index, request := range requests {
		// Step 3.1: Prepare
		pk := new(packet.ItemStackRequest)
		pks = append(pks, pk)

		waiters := make([]chan struct{}, 0)
		handler := newItemStackOperationHandler(
			api.Container(),
			newVirtualInventories(api.Inventories()),
			newResponseMapping(),
		)

		if len(request) == 0 {
			continue
		}

		// Step 3.2: If can inline
		if request[0].CanInline() {
			requestID := api.ItemStackOperation().NewRequestID()
			actions := make([]protocol.StackRequestAction, 0)

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
				}

				if err != nil {
					return false, nil, nil, fmt.Errorf("Commit: %v", err)
				}

				actions = append(actions, result...)
			}

			pk.Requests = append(
				pk.Requests,
				protocol.ItemStackRequest{
					RequestID: int32(requestID),
					Actions:   actions,
				},
			)

			idx := index
			channel := make(chan struct{})
			waiters = append(waiters, channel)

			api.ItemStackOperation().AddNewRequest(
				requestID,
				handler.responseMapping.mapping,
				nil,
				func(response *protocol.ItemStackResponse) {
					mu.Lock()
					defer mu.Unlock()
					serverResponse[idx] = append(serverResponse[idx], response)
					close(channel)
				},
			)
		}

		// Step 3.2: If can not inline
		if !request[0].CanInline() {
			for _, operation := range request {
				var (
					itemNewName *string
					updater     map[resources_control.SlotLocation]resources_control.ExpectedNewItem

					result []protocol.StackRequestAction
					err    error

					requestID resources_control.ItemStackRequestID = api.ItemStackOperation().NewRequestID()
				)

				switch op := operation.(type) {
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

				newRequest := protocol.ItemStackRequest{
					RequestID: int32(requestID),
					Actions:   result,
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
						mu.Lock()
						defer mu.Unlock()
						serverResponse[idx] = append(serverResponse[idx], response)
						close(channel)
					},
				)
			}
		}

		// Step 3.3: Send packet
		err = api.WritePacket(pk)
		if err != nil {
			return false, nil, nil, fmt.Errorf("Commit: %v", err)
		}

		// Step 3.4: Wait changes
		for _, waiter := range waiters {
			<-waiter
		}
	}

	// Setp 4: Return unsuccess
	for _, responses := range serverResponse {
		for _, response := range responses {
			if response.Status != protocol.ItemStackResponseStatusOK {
				return false, pks, serverResponse, nil
			}
		}
	}

	// Step 5: Return success
	i.Discord()
	return true, pks, serverResponse, nil
}
