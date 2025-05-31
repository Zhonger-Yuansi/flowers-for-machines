package resources_control

import "github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"

type (
	// ItemStackRequestID 指示每个物品堆栈操作请求的唯一 ID，
	// 它是以 -1 为首项，公差为 -2 的等差数列
	ItemStackRequestID int32
	// ContainerID 是容器的 ID
	ContainerID uint8
)

// ExpectedNewItem 描述单个物品堆栈在经历一次物品堆栈操作后，
// 其最终应当拥有的一些数据信息。应当说明的是，这些数据信息不
// 会由服务器告知，它应当是客户端内部处理的
type ExpectedNewItem struct {
	// NBTData 指示经过相应的物品堆栈操作后，其 NBT 字段的最终状态。
	// 需要说明的是，物品名称的 NBT 字段无需在此处更改，它会被自动维护
	NBTData map[string]any
}

// ItemStackResponseMapping 是一个由容器 ID 到库存窗口 ID 的映射。
// 由于服务器返回的物品堆栈响应按 ContainerID 来返回更改的物品堆栈，
// 因此本处的资源处理器定义了下面的运行时映射，以便于操作
type ItemStackResponseMapping map[ContainerID]WindowID

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
			item.Stack.NBTData = newData.NBTData
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
