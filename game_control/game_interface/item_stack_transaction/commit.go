package item_stack_transaction

import (
	"fmt"
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface/item_stack_operation"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// Discord 丢弃曾经执行的更改。
// 从本质上说，它将清空底层操作序列
func (i *ItemStackTransaction) Discord() *ItemStackTransaction {
	for index := range i.operations {
		i.operations[index] = nil
	}
	i.operations = nil
	return i
}

// Commit 将底层操作序列内联到单个物品堆栈操作请求数据包中执行物品堆栈操作事务。
// 如果没有返回错误，Commit 在完成后将使用 Discord 清空底层操作序列。
// 应当说明的是，如果事务没有全部成功，则若没有返回错误，则 Discord 仍然会执行。
//
// Commit 在设计上考虑并预期事务的所有都会成功，因此内联将尽可能紧凑，而这依赖于“成功”的预期前提。
// 这意味着，一旦某个步骤失败，那么整个物品堆栈操作都可能失败，并且最终的结果将是未定义的。
//
// success 为真指示该事务的全部操作完全成功，若为否则可能部分失败。
// 作为一种特殊情况，如果底层操作序列为空，则 success 总是真。
//
// pk 指示最终编译得到的数据包，它可以用于调试，但不应重新用于发送；
// serverResponse 则指示租赁服针对 pk 中每个物品堆栈操作的结果
func (i *ItemStackTransaction) Commit() (
	success bool,
	pk *packet.ItemStackRequest,
	serverResponse []*protocol.ItemStackResponse,
	err error,
) {
	if len(i.operations) == 0 {
		return true, nil, make([]*protocol.ItemStackResponse, 0), nil
	}

	api := i.api
	mu := new(sync.Mutex)

	pk = new(packet.ItemStackRequest)
	requests := make([][]item_stack_operation.ItemStackOperation, 0)
	waiters := make([]chan struct{}, 0)

	handler := newItemStackOperationHandler(
		api.Container(),
		api.ConstantPacket(),
		newVirtualInventories(api.Inventories()),
		newResponseMapping(),
	)

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
	serverResponse = make([]*protocol.ItemStackResponse, len(requests))

	// Step 2: Construct actions
	for index, request := range requests {
		if len(request) == 0 {
			continue
		}

		// Step 2.1: If can inline
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
					serverResponse[idx] = response
					close(channel)
				},
			)
		}

		// Step 2.2: If can not inline
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
					newItem := i.api.ConstantPacket().CreativeItemByCNI(op.CreativeItemNetworkID)
					updater = make(map[resources_control.SlotLocation]resources_control.ExpectedNewItem)
					updater[op.Path] = resources_control.ExpectedNewItem{
						NetworkID: newItem.Item.NetworkID,
						NBTData:   newItem.Item.NBTData,
					}
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
						serverResponse[idx] = response
						close(channel)
					},
				)
			}
		}
	}

	// Step 3: Send packet
	err = api.WritePacket(pk)
	if err != nil {
		return false, nil, nil, fmt.Errorf("Commit: %v", err)
	}

	// Step 4: Wait changes
	for _, waiter := range waiters {
		<-waiter
	}

	// Setp 5.1: Return unsuccess
	for _, response := range serverResponse {
		if response.Status != protocol.ItemStackResponseStatusOK {
			_ = i.Discord()
			return false, pk, serverResponse, nil
		}
	}

	// Step 5.2: Return success
	_ = i.Discord()
	return true, pk, serverResponse, nil
}
