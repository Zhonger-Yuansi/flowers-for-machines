package resources_control

import (
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
)

type (
	// SlotID 是单个物品栏槽位的索引，它是从 0 开始索引的
	SlotID uint8
	// Inventory 描述机器人的单个库存
	Inventory struct {
		mu      *sync.RWMutex
		mapping map[SlotID]*protocol.ItemInstance
	}

	// WindowID 是机器人已打开(或持有)的库存的窗口 ID
	WindowID uint32
	// Inventories 描述机器人已打开(或持有)的所有库存，
	// 例如背包、副手和胸甲
	Inventories struct {
		mu      *sync.RWMutex
		mapping map[WindowID]*Inventory
	}

	// SlotLocation 描述一个物品的所在的位置
	SlotLocation struct {
		WindowID WindowID // WindowID 指示该物品所在的库存窗口 ID
		SlotID   SlotID   // SlotID 指示该物品所在库存的槽位索引
	}
)

// NewInventory 返回一个新的 Inventory
func NewInventory() *Inventory {
	return &Inventory{
		mu:      new(sync.RWMutex),
		mapping: make(map[SlotID]*protocol.ItemInstance),
	}
}

// NewInventories 返回一个新的 Inventories
func NewInventories() *Inventories {
	return &Inventories{
		mu:      new(sync.RWMutex),
		mapping: make(map[WindowID]*Inventory),
	}
}

// NewAirItem 返回一个新的空气物品堆栈实例
func NewAirItem() *protocol.ItemInstance {
	return &protocol.ItemInstance{
		StackNetworkID: 0,
		Stack: protocol.ItemStack{
			ItemType: protocol.ItemType{
				NetworkID:     0,
				MetadataValue: 0,
			},
			BlockRuntimeID: 0,
			Count:          0,
			NBTData:        map[string]interface{}(nil),
			CanBePlacedOn:  []string(nil),
			CanBreak:       []string(nil),
			HasNetworkID:   false,
		},
	}
}

// GetItemStack 返回当前库存中物品栏编号为 slotID 的物品堆栈信息。
// 如果不存在，确保返回一个新的空气物品的堆栈实例表示，而非空指针
func (i *Inventory) GetItemStack(slotID SlotID) *protocol.ItemInstance {
	i.mu.RLock()
	defer i.mu.RUnlock()

	result, ok := i.mapping[slotID]
	if !ok {
		return NewAirItem()
	}

	return result
}

// setItemStack 将 item 所指示的物品堆栈实例储存到当前库存的 slotID 处。
// 如果 item 为空指针，则储存为空气；
// 如果 item 未更改且 slotID 处已存在物品，则不作额外操作。
//
// setItemStack 是一个内部实现细节，不应被其他人所使用
func (i *Inventory) setItemStack(slotID SlotID, item *protocol.ItemInstance) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if item == nil {
		i.mapping[slotID] = NewAirItem()
		return
	}

	if item.Stack.NetworkID == -1 {
		if _, ok := i.mapping[slotID]; !ok {
			i.mapping[slotID] = NewAirItem()
		}
		return
	}

	i.mapping[slotID] = item
}

// GetInventory 返回窗口 ID 为 windowID 的库存。
// 如果目标库存不存在，则返回的 existed 为假
func (i *Inventories) GetInventory(windowID WindowID) (inventory *Inventory, existed bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	inventory, existed = i.mapping[windowID]
	return
}

// createInventory 创建一个窗口 ID 为 windowID 的库存。
// 如果库存已经存在，则不会进行任何操作。
//
// createInventory 是一个内部实现细节，不应被其他人所使用
func (i *Inventories) createInventory(windowID WindowID) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.mapping[windowID]; !ok {
		i.mapping[windowID] = NewInventory()
	}
}

// deleteInventory 将窗口 ID 为 windowID 的库存从底层删除。
// 如果库存本身不存在，则不会进行任何操作。
//
// deleteInventory 是一个内部实现细节，不应被其他人所使用
func (i *Inventories) deleteInventory(windowID WindowID) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.mapping[windowID]; ok {
		delete(i.mapping, windowID)
		newMapping := make(map[WindowID]*Inventory)
		for key, value := range i.mapping {
			newMapping[key] = value
		}
		i.mapping = newMapping
	}
}

// GetItemStack 加载位于 windowID 的库存中索引为 slotID 的物品。
// 如果目标库存不存在，则返回的 inventoryExisted 为假
func (i *Inventories) GetItemStack(windowID WindowID, slotID SlotID) (item *protocol.ItemInstance, inventoryExisted bool) {
	inventory, existed := i.GetInventory(windowID)
	if !existed {
		return nil, false
	}
	return inventory.GetItemStack(slotID), true
}

// setItemStack 设置位于 windowsID 库存中索引为 slotID 的物品的数据为 item。
//
// 如果窗口 ID 为 windowID 的库存不存在，则尝试创建其；
// 如果 item 为空指针，则设置为空气；
// 如果 item 未更改且 slotID 处已存在物品，则不作额外操作。
//
// setItemStack 是一个内部实现细节，不应被其他人所使用
func (i *Inventories) setItemStack(windowID WindowID, slotID SlotID, item *protocol.ItemInstance) {
	for {
		i.createInventory(windowID)

		inventory, existed := i.GetInventory(windowID)
		if !existed {
			continue
		}

		inventory.setItemStack(slotID, item)
		break
	}
}
