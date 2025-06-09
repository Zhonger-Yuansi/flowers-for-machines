package nbt_parser

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
)

type BookNBT struct {
	Pages  []string
	Author string
	Title  string
}

type Book struct {
	DefaultItem
	NBT BookNBT
}

// parse ..
func (b *Book) parse(tag map[string]any) {
	b.DefaultItem.Enhance.ItemComponent.LockInInventory = false
	b.DefaultItem.Enhance.ItemComponent.LockInSlot = false
	b.DefaultItem.Enhance.EnchList = nil
	b.DefaultItem.Block = ItemBlockData{}

	if len(tag) == 0 {
		return
	}

	pages, _ := tag["pages"].([]any)
	for _, page := range pages {
		content, ok := page.(string)
		if !ok {
			continue
		}
		b.NBT.Pages = append(b.NBT.Pages, content)
	}

	b.NBT.Author, _ = tag["author"].(string)
	b.NBT.Title, _ = tag["title"].(string)
}

func (b *Book) ParseNormal(nbtMap map[string]any) error {
	err := b.DefaultItem.ParseNormal(nbtMap)
	if err != nil {
		return fmt.Errorf("ParseNormal: %v", err)
	}

	tag, _ := nbtMap["tag"].(map[string]any)
	b.parse(tag)

	return nil
}

func (b *Book) ParseNetwork(item protocol.ItemStack, itemNetworkIDToName map[int32]string) error {
	err := b.DefaultItem.ParseNetwork(item, itemNetworkIDToName)
	if err != nil {
		return fmt.Errorf("ParseNetwork: %v", err)
	}
	b.parse(item.NBTData)
	return nil
}

func (b *Book) NeedSpecialHandle() bool {
	if b.ItemName() == "minecraft:written_book" {
		return true
	}
	for _, page := range b.NBT.Pages {
		if len(page) > 0 {
			return true
		}
	}
	return false
}

func (d Book) NeedCheckCompletely() bool {
	return true
}

func (b *Book) TypeStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	basicInfo := b.DefaultItem.TypeStableBytes()
	w.ByteSlice(&basicInfo)
	protocol.FuncSliceUint16Length(w, &b.NBT.Pages, w.String)
	w.String(&b.NBT.Author)
	w.String(&b.NBT.Title)

	return buf.Bytes()
}

func (b *Book) FullStableBytes() []byte {
	return append(b.TypeStableBytes(), b.Basic.Count)
}
