package game_interface

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface/item_stack_operation"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
)

// ------------------------- ItemStackOperation -------------------------

// ItemStackOperation 是物品操作请求的包装实现
type ItemStackOperation struct {
	api *resources_control.Resources
}

// NewItemStackOperation 基于 api 创建并返回一个新的 ItemStackOperation
func NewItemStackOperation(api *resources_control.Resources) *ItemStackOperation {
	return &ItemStackOperation{api: api}
}

// NewTransaction 打开一个新的物品堆栈操作事务。
// 多个事务可以被同时打开，但各个事务的操作内容
// 不应该发生重叠，否则结果是未知的。
func (i *ItemStackOperation) NewTransaction() *ItemStackTransaction {
	return &ItemStackTransaction{
		api:        i.api,
		operations: nil,
	}
}

// ------------------------- ItemStackTransaction -------------------------

// ItemStackTransaction 是单个物品操作事务，
// 它希望使用者尽可能多的将物品堆栈请求内联在一个数据包中，
// 这样可以有效的节省操作的时间消耗
type ItemStackTransaction struct {
	api            *resources_control.Resources
	operations     []item_stack_operation.ItemStackOperation
	stackNetworkID map[resources_control.SlotLocation]int32
}

// ------------------------- ItemStackTransaction & Commit Operations -------------------------

// slotLocationToContainerID 根据 slotLocation查找本地持有(或打开) 的库存，
// 并找到对应 slotLocation 的容器 ID
func (i *ItemStackTransaction) slotLocationToContainerID(slotLocation resources_control.SlotLocation) (
	result resources_control.ContainerID,
	found bool,
) {
	switch slotLocation.WindowID {
	case protocol.WindowIDInventory:
		return protocol.ContainerCombinedHotBarAndInventory, true
	case protocol.WindowIDOffHand:
		return protocol.ContainerOffhand, true
	case protocol.WindowIDArmour:
		return protocol.ContainerArmor, true
	case protocol.WindowIDUI:
		return 0, false // TODO: Figure out what WindowIDUI means
	}

	containerData, existed := i.api.Container().ContainerData()
	if !existed {
		return 0, false
	}
	if containerData.WindowID != byte(slotLocation.WindowID) {
		return 0, false
	}

	containerTypeWithSlot := mapping.ContainerTypeWithSlot{
		ContainerType: int(containerData.ContainerType),
	}
	if mapping.ContainerNeedSlotIDMapping[containerTypeWithSlot.ContainerType] {
		containerTypeWithSlot.SlotID = uint8(slotLocation.SlotID)
	}
	if result, ok := mapping.ContainerIDMapping[containerTypeWithSlot]; ok {
		if result == mapping.ContainerIDUnknown || result == mapping.ContainerIDCanNotOpen {
			return 0, false
		}
		return resources_control.ContainerID(result), true
	}

	return 0, false
}

// Discord 丢弃曾经执行的更改。
// 从本质上说，它将清空底层操作序列
func (i *ItemStackTransaction) Discord() {
	i.operations = nil
}

// // Commit 将底层操作序列内联到单个的物品堆栈请求的数据包中，并将它发送至租赁服。
// // Commit 在设计上考虑并预期整个事务都会成功。
// func (i *ItemStackTransaction) Commit() (success bool, serverResponse []protocol.ItemStackResponse, err error) {
// 	pk := packet.ItemStackRequest{}

// }

// ------------------------- ItemStackTransaction & Add Operations -------------------------

// MoveItem 将 source 处的物品移动到 destination 处，
// 且只移动 count 个物品
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
// 此操作需要保证背包已被打开，或者已打开的容器中可以在背包中移动物品
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

// SwapItem 交换 source 处和 destination 处的物品
func (i *ItemStackTransaction) SwapItem(source resources_control.SlotLocation, destination resources_control.SlotLocation) {
	i.operations = append(i.operations, item_stack_operation.Swap{
		Source:      source,
		Destination: destination,
	})
}

// DropItem 将 slot 处的物品丢出，且只丢出 count 个
func (i *ItemStackTransaction) DropItem(slot resources_control.SlotLocation, count uint8) {
	i.operations = append(i.operations, item_stack_operation.Drop{
		Path:  slot,
		Count: count,
	})
}

// DropItem 将背包中 slot 处的物品丢出，且只丢出 count 个。
// 此操作需要保证背包已被打开，或者已打开的容器中可以在背包中移动物品
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
// 不同于 DropInventoryItem，DropHotbarItem 可以在未打开背包时使用
// TODO: 需要验证
func (i *ItemStackTransaction) DropHotbarItem(slot resources_control.SlotID, count uint8) {
	i.operations = append(i.operations, item_stack_operation.DropHotbar{
		SlotID: slot,
		Count:  count,
	})
}

// GetCreativeItem 从创造物品栏获取 创造物品网络 ID 为
// creativeItemNetworkID 的物品到背包中的 slot 处，
// 且只移动 count 个物品
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
// 如果操作成功，则物品将回到原位
func (i *ItemStackTransaction) RenameItem(slot resources_control.SlotLocation, newName string) {
	i.operations = append(i.operations, item_stack_operation.Renaming{
		Path:    slot,
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
// 如果操作成功，则新旗帜将回到原位
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
