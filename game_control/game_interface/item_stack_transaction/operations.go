package item_stack_transaction

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface/item_stack_operation"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// MoveItem 将 source 处的物品移动到 destination 处，
// 且只移动 count 个物品。
//
// 该操作是支持内联的，它会与所有支持内联的操作一起被
// 内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) MoveItem(
	source resources_control.SlotLocation,
	destination resources_control.SlotLocation,
	count uint8,
) {
	i.operations = append(i.operations, item_stack_operation.Move{
		Source:      source,
		Destination: destination,
		Count:       int32(count),
	})
}

// MoveInventoryItem 将背包中 source 处的物品移动到 destination 处，
// 且只移动 count 个物品。
//
// 此操作需要保证背包已被打开，或者已打开的容器中可以在背包中移动物品。
//
// 该操作是支持内联的，它会与所有支持内联的操作一起被
// 内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) MoveInventoryItem(
	source resources_control.SlotID,
	destination resources_control.SlotID,
	count uint8,
) {
	i.operations = append(i.operations, item_stack_operation.Move{
		Source: resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   source,
		},
		Destination: resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   destination,
		},
		Count: int32(count),
	})
}

// SwapItem 交换 source 处和 destination 处的物品。
//
// 该操作是支持内联的，它会与所有支持内联的操作一起被
// 内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) SwapItem(source resources_control.SlotLocation, destination resources_control.SlotLocation) {
	i.operations = append(i.operations, item_stack_operation.Swap{
		Source:      source,
		Destination: destination,
	})
}

// DropItem 将 slot 处的物品丢出，且只丢出 count 个。
//
// 该操作是支持内联的，它会与所有支持内联的操作一起被
// 内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) DropItem(slot resources_control.SlotLocation, count uint8) {
	i.operations = append(i.operations, item_stack_operation.Drop{
		Path:  slot,
		Count: count,
	})
}

// DropItem 将背包中 slot 处的物品丢出，
// 且只丢出 count 个。
//
// 此操作需要保证背包已被打开，
// 或者已打开的容器中可以在背包中移动物品。
//
// 该操作是支持内联的，它会与所有支持内联
// 的操作一起被内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) DropInventoryItem(slot resources_control.SlotID, count uint8) {
	i.operations = append(i.operations, item_stack_operation.Drop{
		Path: resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   slot,
		},
		Count: count,
	})
}

// DropItem 将背包中 slot 处的物品丢出，且只丢出 count 个。
// 不同于 DropInventoryItem，DropHotbarItem 可以在未打开背包时使用。
//
// 该操作是支持内联的，它会与所有支持内联的操作一起被
// 内联到单个物品堆栈操作请求中
func (i *ItemStackTransaction) DropHotbarItem(slot resources_control.SlotID, count uint8) {
	i.operations = append(i.operations, item_stack_operation.DropHotbar{
		SlotID: slot,
		Count:  count,
	})
}

// GetCreativeItem 从创造物品栏获取 创造物品网络 ID 为
// creativeItemNetworkID 的物品到背包中的 slot 处，
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
) {
	i.operations = append(i.operations, item_stack_operation.CreativeItem{
		CreativeItemNetworkID: creativeItemNetworkID,
		SlotID:                slot.SlotID,
		Count:                 count,
	})
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
// 请确保不要将此操作与内联操作混用，否则将会发送多个
// 数据包从而降低事务的效率。
//
// 除此外，基于非内联操作的并发组织，你无法在同一个物品
// 栏处重用非内联操作 (重命名操作或织布机操作)，除非您在
// 操作前引入了至少一个内联操作，否则整个事务将会失败
func (i *ItemStackTransaction) RenameItem(slot resources_control.SlotLocation, newName string) {
	i.operations = append(i.operations, item_stack_operation.Renaming{
		Path:    slot,
		NewName: newName,
	})
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
func (i *ItemStackTransaction) RenameInventoryItem(slot resources_control.SlotID, newName string) {
	i.operations = append(i.operations, item_stack_operation.Renaming{
		Path: resources_control.SlotLocation{
			WindowID: protocol.WindowIDInventory,
			SlotID:   slot,
		},
		NewName: newName,
	})
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
) {
	i.operations = append(i.operations, item_stack_operation.Looming{
		UsePattern:  len(patternName) > 0,
		PatternName: patternName,
		PatternPath: patternSlot,
		BannerPath:  bannerSlot,
		DyePath:     dyeSlot,
		ResultItem:  resultItem,
	})
}
