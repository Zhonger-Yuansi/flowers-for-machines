package nbt_assigner_interface

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
)

var (
	NBTItemIsSupported    func(item nbt_parser_interface.Item) bool
	MakeNBTItemMethod     func(console *nbt_console.Console, cache *nbt_cache.NBTCacheSystem, multipleItems ...nbt_parser_interface.Item) []Item
	EnchMultiple          func(console *nbt_console.Console, multipleItems [27]*nbt_parser_interface.Item) error
	RenameMultiple        func(console *nbt_console.Console, multipleItems [27]*nbt_parser_interface.Item) error
	EnchAndRenameMultiple func(console *nbt_console.Console, multipleItems [27]*nbt_parser_interface.Item) error
)

type Item interface {
	Append(item ...nbt_parser_interface.Item)
	Make() (resultSlot map[uint64]resources_control.SlotID, err error)
}
