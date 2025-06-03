package item_stack_transaction

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface/item_stack_operation"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// MoveItem 将 source 处的物品移动到 destination 处，
// 且只移动 count 个物品。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一
// 起被内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) MoveItem(
	source resources_control.SlotLocation,
	destination resources_control.SlotLocation,
	count uint8,
) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.Move{
		Source:      source,
		Destination: destination,
		Count:       int32(count),
	})
	return i
}

// MoveBetweenInventory 将背包中 source 处的物品移动到 destination 处，
// 且只移动 count 个物品。
//
// 此操作需要保证背包已被打开，或者已打开的容器中可以在背包中移动物品。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一 起被内联到单个
// 物品堆栈操作请求中
func (i *ItemStackTransaction) MoveBetweenInventory(
	source resources_control.SlotID,
	destination resources_control.SlotID,
	count uint8,
) *ItemStackTransaction {
	return i.MoveItem(
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   source,
		},
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   destination,
		},
		count,
	)
}

// MoveToContainer 将背包中 source 处的物品移动到已打开容器的
// destination 处，且只移动 count 个物品。
//
// 此操作需要保证目前已经打开了一个容器，否则效果将会与
// MoveInventoryItem 等同。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一起被内联
// 到单个物品堆栈操作请求中
func (i *ItemStackTransaction) MoveToContainer(
	source resources_control.SlotID,
	destination resources_control.SlotID,
	count uint8,
) *ItemStackTransaction {
	data, _ := i.api.Container().ContainerData()
	return i.MoveItem(
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   source,
		},
		resources_control.SlotLocation{
			WindowID: resources_control.WindowID(data.WindowID),
			SlotID:   destination,
		},
		count,
	)
}

// MoveToInventory 将已打开容器中 source 处的物品移动到
// 背包的 destination 处，且只移动 count 个物品。
//
// 此操作需要保证目前已经打开了一个容器，否则效果将会与
// MoveInventoryItem 等同。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一起
// 被内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) MoveToInventory(
	source resources_control.SlotID,
	destination resources_control.SlotID,
	count uint8,
) *ItemStackTransaction {
	data, _ := i.api.Container().ContainerData()
	return i.MoveItem(
		resources_control.SlotLocation{
			WindowID: resources_control.WindowID(data.WindowID),
			SlotID:   source,
		},
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   destination,
		},
		count,
	)
}

// SwapItem 交换 source 处和 destination 处的物品。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作
// 一起被内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) SwapItem(
	source resources_control.SlotLocation,
	destination resources_control.SlotLocation,
) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.Swap{
		Source:      source,
		Destination: destination,
	})
	return i
}

// SwapBetweenInventory 交换背包中 source
// 处和背包中 destination 处的物品。
//
// 此操作需要保证背包已被打开，或者已打开
// 的容器中可以在背包中移动物品。
//
// 该操作是支持内联的，它会与所有相邻的支
// 持内联的操作一起被内联到单个物品堆栈操
// 作请求中
func (i *ItemStackTransaction) SwapBetweenInventory(
	source resources_control.SlotID,
	destination resources_control.SlotID,
) *ItemStackTransaction {
	return i.SwapItem(
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   source,
		},
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   destination,
		},
	)
}

// SwapInventoryBetweenContainer 交换背包中 source
// 处和已打开容器 destination 处的物品。
//
// 此操作需要保证目前已经打开了一个容器，否则效果将会与
// SwapInventoryItem 等同。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一起
// 被内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) SwapInventoryBetweenContainer(
	source resources_control.SlotID,
	destination resources_control.SlotID,
) *ItemStackTransaction {
	data, _ := i.api.Container().ContainerData()
	return i.SwapItem(
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   source,
		},
		resources_control.SlotLocation{
			WindowID: resources_control.WindowID(data.WindowID),
			SlotID:   destination,
		},
	)
}

// DropItem 将 slot 处的物品丢出，且只丢出 count 个。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一
// 起被内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) DropItem(slot resources_control.SlotLocation, count uint8) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.Drop{
		Path:  slot,
		Count: count,
	})
	return i
}

// DropItem 将背包中 slot 处的物品丢出，
// 且只丢出 count 个。
//
// 此操作需要保证背包已被打开，
// 或者已打开的容器中可以在背包中移动物品。
//
// 该操作是支持内联的，它会与所有相邻的支
// 持内联的操作一起被内联到单个物品堆栈操
// 作请求中
func (i *ItemStackTransaction) DropInventoryItem(slot resources_control.SlotID, count uint8) *ItemStackTransaction {
	return i.DropItem(
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   slot,
		},
		count,
	)
}

// DropItem 将快捷栏 slot 处的物品丢出，且只丢出 count 个。
//
// DropHotbarItem 与 DropInventoryItem 不同之处在于其可以
// 在未打开背包时使用。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一起被
// 内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) DropHotbarItem(slot resources_control.SlotID, count uint8) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.DropHotbar{
		SlotID: slot,
		Count:  count,
	})
	return i
}

// DropItem 将已打开容器 slot 处的物品丢出，且只丢出 count 个。
//
// 此操作需要保证目前已经打开了一个容器，否则效果将会与
// DropInventoryItem 等同。
//
// 该操作是支持内联的，它会与所有相邻的支持内联的操作一起被内联
// 到单个物品堆栈操作请求中
func (i *ItemStackTransaction) DropContainerItem(slot resources_control.SlotID, count uint8) *ItemStackTransaction {
	data, _ := i.api.Container().ContainerData()
	return i.DropItem(
		resources_control.SlotLocation{
			WindowID: resources_control.WindowID(data.WindowID),
			SlotID:   slot,
		},
		count,
	)
}

// GetCreativeItem 从创造物品栏获取 创造物品网络 ID 为
// creativeItemNetworkID 的物品到 slot 处，
// 且只移动 count 个物品。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作混用，否则将会发送多个
// 数据包从而降低事务的效率
func (i *ItemStackTransaction) GetCreativeItem(
	creativeItemNetworkID uint32,
	slot resources_control.SlotLocation,
	count uint8,
) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.CreativeItem{
		UseCreativeItemNetworkID: true,
		CreativeItemNetworkID:    creativeItemNetworkID,
		UseNetworkID:             false,
		NetworkID:                0,
		Path:                     slot,
		Count:                    count,
	})
	return i
}

// GetCreativeItemToInventory 从创造物品栏获取 创造物品网络 ID 为
// creativeItemNetworkID 的物品到背包中的 slot 处，
// 且只移动 count 个物品。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作混用，否则将会发送多个
// 数据包从而降低事务的效率
func (i *ItemStackTransaction) GetCreativeItemToInventory(
	creativeItemNetworkID uint32,
	slot resources_control.SlotID,
	count uint8,
) *ItemStackTransaction {
	return i.GetCreativeItem(
		creativeItemNetworkID,
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   slot,
		},
		count,
	)
}

// GetCreativeItemByNetworkID 从创造物品栏获取
// 网络数字 ID 为 networkID 的物品到 slot 处，
// 且只移动 count 个物品。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作混用，否则将会发送多个
// 数据包从而降低事务的效率
func (i *ItemStackTransaction) GetCreativeItemByNetworkID(
	networkID int32,
	slot resources_control.SlotLocation,
	count uint8,
) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.CreativeItem{
		UseCreativeItemNetworkID: false,
		CreativeItemNetworkID:    0,
		UseNetworkID:             true,
		NetworkID:                networkID,
		Path:                     slot,
		Count:                    count,
	})
	return i
}

// GetCreativeItemToInventoryByNetworkID 从创造物品栏获取
// 网络数字 ID 为 networkID 的物品到背包中的 slot 处，且只移
// 动 count 个物品。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作混用，否则将会发送多个
// 数据包从而降低事务的效率
func (i *ItemStackTransaction) GetCreativeItemToInventoryByNetworkID(
	networkID int32,
	slot resources_control.SlotID,
	count uint8,
) *ItemStackTransaction {
	return i.GetCreativeItemByNetworkID(
		networkID,
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   slot,
		},
		count,
	)
}

// RenameItem 将 slot 处的物品全部重命名为 newName。
//
// 重命名操作是通过铁砧完成的，这意味着您需要确保铁砧已被打开，
// 且铁砧内没有放置任何物品。
//
// 如果操作成功，则物品将回到原位。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作交替使用，而是应该尽可能
// 连续的使用多个非内联操作。如果不这么做，提交事务时将
// 会发送多个数据包从而降低事务的效率。
//
// 除此外，基于非内联操作的并发组织，你无法在同一个物品
// 栏处重用非内联操作 (重命名操作或织布机操作)，除非您在
// 操作前引入了至少一个内联操作，否则整个事务将会失败
func (i *ItemStackTransaction) RenameItem(slot resources_control.SlotLocation, newName string) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.Renaming{
		Path:    slot,
		NewName: newName,
	})
	return i
}

// RenameInventoryItem 将背包中 slot 处的物品全部重命名为 newName。
//
// 重命名操作是通过铁砧完成的，这意味着您需要确保铁砧已被打开，
// 且铁砧内没有放置任何物品。如果操作成功，则物品将回到原位。
//
// 与 RenameItem 的不同之处在于，它只能操作背包中的物品，
// 因此您需要确保背包已被打开。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作交替使用，而是应该尽可能
// 连续的使用多个非内联操作。如果不这么做，提交事务时将
// 会发送多个数据包从而降低事务的效率。
//
// 除此外，基于非内联操作的并发组织，你无法在同一个物品
// 栏处重用非内联操作 (重命名操作或织布机操作)，除非您在
// 操作前引入了至少一个内联操作，否则整个事务将会失败
func (i *ItemStackTransaction) RenameInventoryItem(
	slot resources_control.SlotID,
	newName string,
) *ItemStackTransaction {
	return i.RenameItem(
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   slot,
		},
		newName,
	)
}

// Looming 将 patternSlot 处的旗帜放入织布机中，
// 并通过使用 dyeSlot 处的染料合成新旗帜。
//
// patternName 是织布时使用的图案，patternSlot
// 则指示该图案物品的位置。如果无需使用图案，
// 请将 patternName 和 patternSlot 都置为默认的
// 零值。
//
// resultItem 指示期望得到的旗帜的部分数据。
// 如果操作成功，则新旗帜将回到原位。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作交替使用，而是应该尽可能
// 连续的使用多个非内联操作。如果不这么做，提交事务时将
// 会发送多个数据包从而降低事务的效率。
//
// 除此外，基于非内联操作的并发组织，你无法在同一个物品
// 栏处重用非内联操作 (重命名操作或织布机操作)，除非您在
// 操作前引入了至少一个内联操作，否则整个事务将会失败
func (i *ItemStackTransaction) Looming(
	patternName string,
	patternSlot resources_control.SlotLocation,
	bannerSlot resources_control.SlotLocation,
	dyeSlot resources_control.SlotLocation,
	resultItem resources_control.ExpectedNewItem,
) *ItemStackTransaction {
	i.operations = append(i.operations, item_stack_operation.Looming{
		UsePattern:  len(patternName) > 0,
		PatternName: patternName,
		PatternPath: patternSlot,
		BannerPath:  bannerSlot,
		DyePath:     dyeSlot,
		ResultItem:  resultItem,
	})
	return i
}

// LoomingFromInventory 将背包中 patternSlot 处的旗帜放入织布机中，
// 并通过使用背包中 dyeSlot 处的染料合成新旗帜。
//
// patternName 是织布时使用的图案，patternSlot 则指示该图案物品
// 在背包中的位置。如果无需使用图案，请将 patternName 和 patternSlot
// 都置为默认的零值。
//
// resultItem 指示期望得到的旗帜的部分数据。
// 如果操作成功，则新旗帜将回到原位。
//
// 该操作不支持内联，但任何不支持内联的操作都可以被并发，
// 这意味着虽然它们会被分配在各自独立的物品堆栈请求中，
// 但最终可以被紧缩在一个数据包中。
//
// 请确保不要将此操作与内联操作交替使用，而是应该尽可能
// 连续的使用多个非内联操作。如果不这么做，提交事务时将
// 会发送多个数据包从而降低事务的效率。
//
// 除此外，基于非内联操作的并发组织，你无法在同一个物品
// 栏处重用非内联操作 (重命名操作或织布机操作)，除非您在
// 操作前引入了至少一个内联操作，否则整个事务将会失败
func (i *ItemStackTransaction) LoomingFromInventory(
	patternName string,
	patternSlot resources_control.SlotID,
	bannerSlot resources_control.SlotID,
	dyeSlot resources_control.SlotID,
	resultItem resources_control.ExpectedNewItem,
) *ItemStackTransaction {
	return i.Looming(
		patternName,
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   patternSlot,
		},
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   bannerSlot,
		},
		resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   dyeSlot,
		},
		resultItem,
	)
}
