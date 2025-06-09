package item_stack_transaction

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// ------------------------- Define -------------------------

// virtualInventories 是虚拟库存实现，
// 它是一个内部实现细节，不应被其他人使用
type virtualInventories struct {
	api            *resources_control.Inventories
	stackNetworkID map[resources_control.SlotLocation]int32
	itemCount      map[resources_control.SlotLocation]uint8
}

// newVirtualInventories 基于 api 创建一个新的 newVirtualInventories
func newVirtualInventories(api *resources_control.Inventories) *virtualInventories {
	return &virtualInventories{
		api:            api,
		stackNetworkID: make(map[resources_control.SlotLocation]int32),
		itemCount:      make(map[resources_control.SlotLocation]uint8),
	}
}

// ------------------------- Stack Network ID -------------------------

// allStackNetworkID 列出窗口 ID 为 windowID 的库存中所有物品的运行时 ID
func (v *virtualInventories) allStackNetworkID(windowID resources_control.WindowID) map[resources_control.SlotID]int32 {
	result := make(map[resources_control.SlotID]int32)

	mapping, inventoryExisted := v.api.GetAllItemStack(windowID)
	if !inventoryExisted {
		return nil
	}
	for slotID, item := range mapping {
		result[slotID] = item.StackNetworkID
	}

	for location, stackNetworkID := range v.stackNetworkID {
		if location.WindowID != windowID {
			continue
		}
		result[location.SlotID] = stackNetworkID
	}

	return result
}

// loadStackNetworkID 加载 slotLocation 处的物品堆栈网络 ID
func (v *virtualInventories) loadStackNetworkID(slotLocation resources_control.SlotLocation) (result int32, err error) {
	if result, ok := v.stackNetworkID[slotLocation]; ok {
		return result, nil
	}

	item, inventoryExisted := v.api.GetItemStack(slotLocation.WindowID, slotLocation.SlotID)
	if !inventoryExisted {
		return 0, fmt.Errorf("loadStackNetworkID: Can not find the item whose at %#v", slotLocation)
	}
	v.stackNetworkID[slotLocation] = item.StackNetworkID

	return item.StackNetworkID, nil
}

// setStackNetworkID 设置 slotLocation 处的物品堆栈网络 ID 为 requestID
func (v *virtualInventories) setStackNetworkID(
	slotLocation resources_control.SlotLocation,
	requestID resources_control.ItemStackRequestID,
) {
	v.stackNetworkID[slotLocation] = int32(requestID)
}

// loadAndSetStackNetworkID 加载 slotLocation 处的物品堆栈网络 ID，
// 并将 slotLocation 处的物品堆栈网络 ID 更新为 requestID
func (v *virtualInventories) loadAndSetStackNetworkID(
	slotLocation resources_control.SlotLocation,
	requestID resources_control.ItemStackRequestID,
) (result int32, err error) {
	result, err = v.loadStackNetworkID(slotLocation)
	if err != nil {
		return 0, fmt.Errorf("loadAndSetStackNetworkID: %v", err)
	}
	v.setStackNetworkID(slotLocation, requestID)
	return
}

// ------------------------- Item Count -------------------------

// allItemCount 列出窗口 ID 为 windowID 的库存中所有物品的数量
func (v *virtualInventories) allItemCount(windowID resources_control.WindowID) map[resources_control.SlotID]uint8 {
	result := make(map[resources_control.SlotID]uint8)

	mapping, inventoryExisted := v.api.GetAllItemStack(windowID)
	if !inventoryExisted {
		return nil
	}
	for slotID, item := range mapping {
		result[slotID] = uint8(item.Stack.Count)
	}

	for location, itemCount := range v.itemCount {
		if location.WindowID != windowID {
			continue
		}
		result[location.SlotID] = itemCount
	}

	return result
}

// loadItemCount 加载 slotLocation 处的物品数量
func (v *virtualInventories) loadItemCount(slotLocation resources_control.SlotLocation) (result uint8, err error) {
	if result, ok := v.itemCount[slotLocation]; ok {
		return result, nil
	}

	item, inventoryExisted := v.api.GetItemStack(slotLocation.WindowID, slotLocation.SlotID)
	if !inventoryExisted {
		return 0, fmt.Errorf("loadItemCount: Can not find the item whose at %#v", slotLocation)
	}
	v.itemCount[slotLocation] = uint8(item.Stack.Count)

	return v.itemCount[slotLocation], nil
}

// addItemCount 将 slotLocation 处的物品数量添加 delta。
// 另外，delta 可以是负数。
// allowNoChange 指示是否允许 delta 为 0。如果填写假且
// delta 为 0，那么返回错误
func (v *virtualInventories) addItemCount(
	slotLocation resources_control.SlotLocation,
	delta int8,
	allowNoChange bool,
) error {
	if delta == 0 && !allowNoChange {
		return fmt.Errorf("addItemCount: Item count no change when not allow no change")
	}

	resultCount := int8(v.itemCount[slotLocation]) + delta
	if resultCount < 0 {
		return fmt.Errorf(
			"addItemCount: Invalid result count (origin count = %d, delta = %d, result count = %d)",
			v.itemCount[slotLocation], delta, resultCount,
		)
	}
	v.itemCount[slotLocation] = uint8(resultCount)

	return nil
}

// addItemCount 将 slotLocation 处的物品数量设置为 count
func (v *virtualInventories) setItemCount(slotLocation resources_control.SlotLocation, count uint8) error {
	if count > 64 {
		return fmt.Errorf("setItemCount: Invalid count %d because it more than 64", count)
	}
	v.itemCount[slotLocation] = count
	return nil
}

// loadAndAddItemCount 加载 slotLocation 处的物品数量，
// 并将该数量添加 delta。
// allowNoChange 指示是否允许 delta 为 0。如果填写假且
// delta 为 0，那么返回错误
func (v *virtualInventories) loadAndAddItemCount(
	slotLocation resources_control.SlotLocation,
	delta int8,
	allowNoChange bool,
) (result uint8, err error) {
	result, err = v.loadItemCount(slotLocation)
	if err != nil {
		return 0, fmt.Errorf("loadAndAddItemCount: %v", err)
	}
	err = v.addItemCount(slotLocation, delta, allowNoChange)
	if err != nil {
		return 0, fmt.Errorf("loadAndAddItemCount: %v", err)
	}
	return
}

// setItemCount 加载 slotLocation 处的物品数量，
// 并将该数量设置为 count
func (v *virtualInventories) loadAndSetItemCount(
	slotLocation resources_control.SlotLocation,
	count uint8,
) (result uint8, err error) {
	result, err = v.loadItemCount(slotLocation)
	if err != nil {
		return 0, fmt.Errorf("loadAndSetItemCount: %v", err)
	}
	err = v.setItemCount(slotLocation, count)
	if err != nil {
		return 0, fmt.Errorf("loadAndAddItemCount: %v", err)
	}
	return
}

// ------------------------- End -------------------------
