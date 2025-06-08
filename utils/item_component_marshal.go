package utils

import (
	"encoding/json"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
)

// ItemComponent 是一个物品的物品组件数据
type ItemComponent struct {
	// 控制此物品/方块 (在冒险模式下) 可以使用/放置在其上的方块类型
	CanPlaceOn []string
	// 控制此物品/方块 (在冒险模式下) 可以破坏的方块类型。
	// 此效果不会改变原本的破坏速度和破坏后掉落物
	CanDestroy []string
	// 阻止该物品被从玩家的物品栏
	// 移除、丢弃或用于合成
	LockInInventory bool
	// 阻止该物品被从玩家物品栏的该槽位
	// 移动、移除、丢弃或用于合成
	LockInSlot bool
	// 使该物品在玩家死亡时不会掉落
	KeepOnDeath bool
}

// ParseItemComponent 从 nbtMap 解析一个物品的物品组件数据
func ParseItemComponent(nbtMap map[string]any) (result ItemComponent) {
	list, ok := nbtMap["CanDestroy"].([]any)
	if ok {
		for _, value := range list {
			val, ok := value.(string)
			if !ok {
				continue
			}
			result.CanDestroy = append(result.CanDestroy, val)
		}
	}

	list, ok = nbtMap["CanPlaceOn"].([]any)
	if ok {
		for _, value := range list {
			val, ok := value.(string)
			if !ok {
				continue
			}
			result.CanPlaceOn = append(result.CanPlaceOn, val)
		}
	}

	tag, ok := nbtMap["tag"].(map[string]any)
	if !ok {
		return
	}

	itemLock, _ := tag["minecraft:item_lock"].(byte)
	switch itemLock {
	case 1:
		result.LockInSlot = true
	case 2:
		result.LockInInventory = true
	}

	keepOnDeath, _ := tag["minecraft:keep_on_death"].(byte)
	if keepOnDeath == 1 {
		result.KeepOnDeath = true
	}

	return
}

// ParseItemComponentNetwork 从 item 解析一个物品的物品组件数据
func ParseItemComponentNetwork(item protocol.ItemStack) (result ItemComponent) {
	result.CanDestroy = item.CanBreak
	result.CanPlaceOn = item.CanBePlacedOn

	if item.NBTData == nil {
		return
	}

	itemLock, _ := item.NBTData["minecraft:item_lock"].(byte)
	switch itemLock {
	case 1:
		result.LockInSlot = true
	case 2:
		result.LockInInventory = true
	}

	keepOnDeath, _ := item.NBTData["minecraft:keep_on_death"].(byte)
	if keepOnDeath == 1 {
		result.KeepOnDeath = true
	}

	return
}

// MarshalItemComponent 将 component 序列化为 MC 命令中的物品组件字符串
func MarshalItemComponent(component ItemComponent) string {
	type Blocks struct {
		Blocks []string `json:"blocks"`
	}
	type Mode struct {
		Mode string `json:"mode"`
	}
	type Component struct {
		CanPlaceOn  *Blocks   `json:"can_place_on,omitempty"`
		CanDestroy  *Blocks   `json:"can_destroy,omitempty"`
		ItemLock    *Mode     `json:"item_lock,omitempty"`
		KeepOnDeath *struct{} `json:"keep_on_death,omitempty"`
	}

	c := Component{}
	if len(component.CanDestroy) > 0 {
		c.CanDestroy = &Blocks{Blocks: component.CanDestroy}
	}
	if len(component.CanPlaceOn) > 0 {
		c.CanPlaceOn = &Blocks{Blocks: component.CanPlaceOn}
	}
	if component.LockInInventory {
		c.ItemLock = &Mode{Mode: "lock_in_inventory"}
	}
	if component.LockInSlot {
		c.ItemLock = &Mode{Mode: "lock_in_slot"}
	}
	if component.KeepOnDeath {
		c.KeepOnDeath = &struct{}{}
	}

	result, _ := json.Marshal(c)
	return string(result)
}
