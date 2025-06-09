package nbt_parser_block

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/mapping"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/cespare/xxhash/v2"
)

type ItemWithSlot struct {
	Item nbt_parser_interface.Item
	Slot uint8
}

type ContainerNBT struct {
	Items []ItemWithSlot
}

type Container struct {
	DefaultBlock
	CanOpen    bool
	CustomName string
	NBT        ContainerNBT
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

func (c *Container) Parse(nbtMap map[string]any) error {
	itemList := make([]map[string]any, 0)

	if !mapping.ContainerCanNotOpen[c.BlockName()] {
		c.CanOpen = true
	}
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

		item, err := nbt_parser_interface.ParseNBTItemNormal(value)
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
	return nil
}

func (b Container) StableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	basicInfo := b.DefaultBlock.StableBytes()
	w.ByteSlice(&basicInfo)
	w.String(&b.CustomName)

	itemMapping := make(map[uint8]ItemWithSlot)
	slots := make([]uint8, 0)
	for _, value := range b.NBT.Items {
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

func (b Container) SetBytes() []byte {
	if len(b.NBT.Items) == 0 {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	itemMapping := make(map[uint64]uint16)
	for _, value := range b.NBT.Items {
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
