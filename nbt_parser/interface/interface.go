package nbt_parser_interface

import "github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"

type Block interface {
	BlockName() string
	BlockStates() map[string]any
	BlockStatesString() string
	Parse(nbtMap map[string]any) error
	NeedSpecialHandle() bool
	NeedCheckCompletely() bool
	StableBytes() []byte
}

var ParseBlock func(blockName string, blockStates map[string]any, blockNBT map[string]any) (block Block, err error)

type Item interface {
	ItemName() string
	ItemCount() uint8
	ItemMetadata() int16
	ParseNetwork(item protocol.ItemStack, itemNetworkIDToName map[int32]string) error
	ParseNormal(nbtMap map[string]any) error
	NeedSpecialHandle() bool
	NeedCheckCompletely() bool
	TypeStableBytes() []byte
	FullStableBytes() []byte
}

var (
	ParseNBTItemNormal  func(nbtMap map[string]any) (item Item, err error)
	ParseNBTItemNetwork func(itemStack protocol.ItemStack, itemNetworkIDToName map[int32]string) (item Item, err error)
)
