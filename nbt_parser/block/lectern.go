package nbt_parser_block

import (
	"bytes"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
)

// LecternNBT ..
type LecternNBT struct {
	ItemRotation float32
	HaveBook     bool
	Book         nbt_parser_interface.Item
}

// 讲台
type Lectern struct {
	DefaultBlock
	NBT LecternNBT
}

func (l Lectern) NeedSpecialHandle() bool {
	return l.NBT.HaveBook
}

func (Lectern) NeedCheckCompletely() bool {
	return true
}

func (l *Lectern) Parse(nbtMap map[string]any) error {
	bookMap, ok := nbtMap["book"].(map[string]any)
	if ok {
		book, err := nbt_parser_interface.ParseItemNormal(bookMap)
		if err != nil {
			return fmt.Errorf("Parse: %v", err)
		}
		l.NBT.HaveBook = true
		l.NBT.Book = book
	}
	return nil
}

func (l Lectern) NBTStableBytes() []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	w.Bool(&l.NBT.HaveBook)
	if l.NBT.HaveBook {
		bookStableBytes := l.NBT.Book.TypeStableBytes()
		w.ByteSlice(&bookStableBytes)
	}

	return buf.Bytes()
}

func (l *Lectern) FullStableBytes() []byte {
	return append(l.DefaultBlock.FullStableBytes(), l.NBTStableBytes()...)
}
