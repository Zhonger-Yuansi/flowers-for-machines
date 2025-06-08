package resources_control

import (
	"sync/atomic"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

type (
	// ItemStackRequestID 指示每个物品堆栈操作请求的唯一 ID，
	// 它是以 -1 为首项，公差为 -2 的等差数列
	ItemStackRequestID int32
	// ContainerID 是容器的 ID
	ContainerID uint8

	// ExpectedNewItem 描述单个物品堆栈在经历一次物品堆栈操作后，
	// 其最终应当拥有的一些数据信息。应当说明的是，这些数据信息不
	// 会由服务器告知，它应当是客户端内部处理的
	ExpectedNewItem struct {
		// NetworkID 是该物品的数值网络 ID，它在单个 MC 版本中不会变化。
		// 它可正亦可负，具体取决于其所关注的物品堆栈实例。
		// 另外，为 NetworkID 填写 -1 的值可以保持其不变
		NetworkID int32

		// UseNBTData 指示是否需要采用下方的 NBTData 更新物品的 NBT 数据
		UseNBTData bool
		// NBTData 指示经过相应的物品堆栈操作后，其 NBT 字段的最终状态。
		// 需要说明的是，物品名称的 NBT 字段无需在此处更改，它会被自动维护
		NBTData map[string]any

		// ChangeRepairCost 指示是否需要更新物品的 RepairCost 字段。
		// 应当说明的是，RepairCost 被用于铁砧的惩罚机制
		ChangeRepairCost bool
		// RepairCostDelta 是要修改的 RepairCost 的增量，可以为负
		RepairCostDelta int32
	}

	// ItemStackResponseMapping 是一个由容器 ID 到库存窗口 ID 的映射。
	// 由于服务器返回的物品堆栈响应按 ContainerID 来返回更改的物品堆栈，
	// 因此本处的资源处理器定义了下面的运行时映射，以便于操作
	ItemStackResponseMapping map[ContainerID]WindowID
)

// ItemStackOperationManager 是所有物品堆栈操作的管理者
type ItemStackOperationManager struct {
	// currentItemStackRequestID 是目前物品堆栈请求的累计 RequestID 计数
	currentItemStackRequestID int32
	// itemStackMapping 存放每个物品堆栈操作请求中的 ItemStackResponseMapping
	itemStackMapping utils.SyncMap[ItemStackRequestID, ItemStackResponseMapping]
	// itemStackUpdater 存放每个物品堆栈操作请求中相关物品的更新函数
	itemStackUpdater utils.SyncMap[ItemStackRequestID, map[SlotLocation]ExpectedNewItem]
	// itemStackCallback 存放所有物品堆栈操作请求的回调函数
	itemStackCallback utils.SyncMap[ItemStackRequestID, func(response *protocol.ItemStackResponse)]
}

// NewItemStackOperationManager 创建并返回一个新的 ItemStackOperationManager
func NewItemStackOperationManager() *ItemStackOperationManager {
	return &ItemStackOperationManager{
		currentItemStackRequestID: 1,
	}
}

// NewRequestID 返回一个可以独立使用的新 RequestID
func (i *ItemStackOperationManager) NewRequestID() ItemStackRequestID {
	return ItemStackRequestID(atomic.AddInt32(&i.currentItemStackRequestID, -2))
}

// AddNewRequest 设置一个即将发送的物品堆栈操作请求的钩子函数。
// mapping 是由容器 ID 到库存窗口 ID 的映射；
//
// updater 存放每个物品堆栈操作请求中所涉及的特定物品的更新方式。
// 需要说明的是，它不必为单个物品堆栈请求中所涉及的所有物品都设置 ExpectedNewItem。
// 就目前而言，只有 NBT 会因物品堆栈操作而发生变化的物品需要这么操作。
//
// callback 是收到服务器响应后应该执行的回调函数
func (i *ItemStackOperationManager) AddNewRequest(
	requestID ItemStackRequestID,
	mapping ItemStackResponseMapping,
	updater map[SlotLocation]ExpectedNewItem,
	callback func(response *protocol.ItemStackResponse),
) {
	i.itemStackMapping.Store(requestID, mapping)
	i.itemStackCallback.Store(requestID, callback)
	if len(updater) > 0 {
		i.itemStackUpdater.Store(requestID, updater)
	}
}

// UpdateItem 通过 serverResponse 和 clientExpected 共同评估 item 的新值。
// slotLocation 指示该物品的位置。应当说明的是，相关修改将直接在 item 上进行
func UpdateItem(
	item *protocol.ItemInstance,
	slotLocation SlotLocation,
	serverResponse protocol.StackResponseSlotInfo,
	clientExpected map[SlotLocation]ExpectedNewItem,
) {
	item.Stack.Count = uint16(serverResponse.Count)
	item.StackNetworkID = serverResponse.StackNetworkID

	if clientExpected != nil {
		newData, ok := clientExpected[slotLocation]
		if ok {
			if newData.NetworkID != -1 {
				item.Stack.ItemType.NetworkID = newData.NetworkID
			}
			if newData.UseNBTData {
				item.Stack.NBTData = newData.NBTData
			}
			if newData.ChangeRepairCost {
				if item.Stack.NBTData == nil {
					item.Stack.NBTData = make(map[string]any)
				}
				repairCost, _ := item.Stack.NBTData["RepairCost"].(int32)
				repairCost += newData.RepairCostDelta
				item.Stack.NBTData["RepairCost"] = repairCost
			}
		}
	}

	if len(item.Stack.NBTData) == 0 && len(serverResponse.CustomName) > 0 {
		item.Stack.NBTData = map[string]any{
			"display": map[string]any{
				"Name": serverResponse.CustomName,
			},
		}
	}

	if len(item.Stack.NBTData) > 0 {
		_, displayExisted := item.Stack.NBTData["display"].(map[string]any)

		// 不存在自定义物品名 且 不存在 display
		if len(serverResponse.CustomName) == 0 && !displayExisted {
			return
		}

		// 存在自定义物品名 且 (不)存在 display
		if len(serverResponse.CustomName) > 0 {
			// 存在自定义物品名 且 存在 display
			if displayExisted {
				item.Stack.NBTData["display"].(map[string]any)["Name"] = serverResponse.CustomName
				return
			}
			// 存在自定义物品名 且 不存在 display
			item.Stack.NBTData["display"] = map[string]any{
				"Name": serverResponse.CustomName,
			}
			return
		}

		// 不存在自定义物品名 且 存在 display
		delete(item.Stack.NBTData["display"].(map[string]any), "Name")
		if len(item.Stack.NBTData["display"].(map[string]any)) == 0 {
			delete(item.Stack.NBTData, "display")
		}
	}
}
