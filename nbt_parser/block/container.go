package nbt_parser_block

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/cespare/xxhash/v2"
)

// ItemWithSlot ..
type ItemWithSlot struct {
	Item nbt_parser_interface.Item
	Slot uint8
}

// ContainerNBT ..
type ContainerNBT struct {
	Items         []ItemWithSlot
	ShulkerFacing uint8
}

// 容器
type Container struct {
	DefaultBlock
	CustomName string
	NBT        ContainerNBT
}

// SetShulkerBoxFacing 将 container 的潜影盒朝向设置为 facing。
// SetShulkerBoxFacing 假定 container 可以被断言为 Container。
// 如果不是 Container，则 SetShulkerBoxFacing 将不执行任何操作
func SetShulkerBoxFacing(container nbt_parser_interface.Block, facing uint8) {
	c, ok := container.(*Container)
	if !ok {
		return
	}
	c.NBT.ShulkerFacing = facing
}

func init() {
	nbt_parser_interface.SetShulkerBoxFacing = SetShulkerBoxFacing
}

func (c Container) NeedSpecialHandle() bool {
	if len(c.CustomName) > 0 {
		return true
	}
	if len(c.NBT.Items) > 0 {
		return true
	}
	return false
}

func (c Container) NeedCheckCompletely() bool {
	return true
}

// ConsiderOpenDirection 指示打开目标容器
// 是否需要考虑其打开方向上的障碍物方块，
// 这似乎只对箱子和潜影盒有效
func (c *Container) ConsiderOpenDirection() bool {
	blockName := c.BlockName()
	if strings.Contains(blockName, "chest") || strings.Contains(blockName, "shulker") {
		return true
	}
	return false
}

func (c *Container) Parse(nbtMap map[string]any) error {
	itemList := make([]map[string]any, 0)

	key, ok := mapping.ContainerStorageKey[c.BlockName()]
	if !ok {
		panic("Parse: Should nerver happened")
	}

	iMap, ok := nbtMap[key].(map[string]any)
	if ok {
		itemList = []map[string]any{iMap}
	}
	iList, ok := nbtMap[key].([]any)
	if ok {
		for _, value := range iList {
			val, ok := value.(map[string]any)
			if !ok {
				continue
			}
			itemList = append(itemList, val)
		}
	}

	for _, value := range itemList {
		slotID, _ := value["Slot"].(byte)

		item, err := nbt_parser_interface.ParseItemNormal(value)
		if err != nil {
			return fmt.Errorf("Parse: %v", err)
		}

		c.NBT.Items = append(
			c.NBT.Items,
			ItemWithSlot{
				Item: item,
				Slot: slotID,
			},
		)
	}

	c.CustomName, _ = nbtMap["CustomName"].(string)
	c.NBT.ShulkerFacing, _ = nbtMap["facing"].(byte)
	return nil
}

func (c Container) NBTStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)
	isShulkerBox := strings.Contains(c.BlockName(), "shulker")

	w.String(&c.CustomName)
	w.Bool(&isShulkerBox)
	if isShulkerBox {
		w.Uint8(&c.NBT.ShulkerFacing)
	}

	itemMapping := make(map[uint8]ItemWithSlot)
	slots := make([]uint8, 0)
	for _, value := range c.NBT.Items {
		itemMapping[value.Slot] = value
		slots = append(slots, value.Slot)
	}

	slices.SortStableFunc(slots, func(a uint8, b uint8) int {
		return cmp.Compare(a, b)
	})

	for _, slot := range slots {
		item := itemMapping[slot]
		stableItemBytes := append(item.Item.FullStableBytes(), item.Slot)
		w.ByteSlice(&stableItemBytes)
	}

	return buf.Bytes()
}

func (c *Container) FullStableBytes() []byte {
	return append(c.DefaultBlock.FullStableBytes(), c.NBTStableBytes()...)
}

func (c Container) SetBytes() []byte {
	if len(c.NBT.Items) == 0 {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	itemMapping := make(map[uint64]uint16)
	for _, value := range c.NBT.Items {
		setHashNumber := xxhash.Sum64(value.Item.TypeStableBytes())
		itemMapping[setHashNumber] += uint16(value.Item.ItemCount())
	}

	setHashNumbers := make([]uint64, 0)
	for setHashNumber := range itemMapping {
		setHashNumbers = append(setHashNumbers, setHashNumber)
	}
	slices.SortStableFunc(setHashNumbers, func(a uint64, b uint64) int {
		return cmp.Compare(a, b)
	})

	for _, setHashNumber := range setHashNumbers {
		count := itemMapping[setHashNumber]
		w.Uint64(&setHashNumber)
		w.Uint16(&count)
	}

	return buf.Bytes()
}
