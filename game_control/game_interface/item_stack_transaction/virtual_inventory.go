package item_stack_transaction

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// virtualInventories 是虚拟库存实现，
// 它是一个内部实现细节，不应被其他人使用
type virtualInventories struct {
	api     *resources_control.Inventories
	mapping map[resources_control.SlotLocation]int32
}

// newVirtualInventories 基于 api 创建一个新的 newVirtualInventories
func newVirtualInventories(api *resources_control.Inventories) *virtualInventories {
	return &virtualInventories{
		api:     api,
		mapping: make(map[resources_control.SlotLocation]int32),
	}
}

// loadStackNetworkID 加载 slotLocation 处的物品堆栈网络 ID
func (v *virtualInventories) loadStackNetworkID(slotLocation resources_control.SlotLocation) (result int32, err error) {
	if result, ok := v.mapping[slotLocation]; ok {
		return result, nil
	}

	item, inventoryExisted := v.api.GetItemStack(slotLocation.WindowID, slotLocation.SlotID)
	if !inventoryExisted {
		return 0, fmt.Errorf("loadStackNetworkID: Can not find the item whose at %#v", slotLocation)
	}
	v.mapping[slotLocation] = item.StackNetworkID

	return item.StackNetworkID, nil
}

// setStackNetworkID 设置 slotLocation 处的物品堆栈网络 ID 为 requestID
func (v *virtualInventories) setStackNetworkID(
	slotLocation resources_control.SlotLocation,
	requestID resources_control.ItemStackRequestID,
) {
	v.mapping[slotLocation] = int32(requestID)
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
