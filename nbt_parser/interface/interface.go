package nbt_parser_interface

import "github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"

var (
	// ParseNBTBlock 从方块实体数据 blockNBT 解析一个方块实体。
	// blockName 和 blockStates 分别指示这个方块实体的名称和方块状态
	ParseBlock func(blockName string, blockStates map[string]any, blockNBT map[string]any) (block Block, err error)
	// ParseItemNormal 从 nbtMap 解析一个 NBT 物品。
	// nbtMap 是含有这个物品 tag 标签的父复合标签
	ParseItemNormal func(nbtMap map[string]any) (item Item, err error)
	// ParseItemNetwork 解析网络传输上的物品堆栈实例 item。
	// itemName 是这个物品堆栈实例的名称
	ParseItemNetwork func(itemStack protocol.ItemStack, itemName string) (item Item, err error)
)

// Block 是所有已实现的 NBT 方块的统称
type Block interface {
	// BlockName 返回这个方块的名称
	BlockName() string
	// BlockStates 返回这个方块的方块状态
	BlockStates() map[string]any
	// BlockStatesString 返回这个方块的方块状态的字符串表示
	BlockStatesString() string
	// Parse 从 nbtMap 解析一个方块实体，
	// nbtMap 是这个方块的方块实体数据
	Parse(nbtMap map[string]any) error
	// NeedSpecialHandle 指示在导入这个
	// 方块实体是否需要进行特殊处理。
	// 如果不需要，则方块直接使用命令放置
	NeedSpecialHandle() bool
	// NeedCheckCompletely 指示在完成这个
	// 方块的导入后是否需要检查其完整性。
	// 如果 NeedSpecialHandle 为假，
	// 则 NeedCheckCompletely 不应被使用
	NeedCheckCompletely() bool
	// StableBytes 返回这个方块实体的数据
	// 的稳定唯一表示
	StableBytes() []byte
}

// Item 是所有已实现的 NBT 物品的统称
type Item interface {
	// ItemName 返回这个物品的名称
	ItemName() string
	// ItemCount 返回这个物品的数量
	ItemCount() uint8
	// ItemMetadata 返回这个物品的元数据
	ItemMetadata() int16
	// ParseNetwork 解析网络传输上的物品堆栈实例 item。
	// itemName 是这个物品的名称
	ParseNetwork(item protocol.ItemStack, itemName string) error
	// ParseNormal 从 nbtMap 解析一个 NBT 物品。
	// nbtMap 是含有这个物品 tag 标签的父复合标签
	ParseNormal(nbtMap map[string]any) error
	// UnderlyingItem 返回这个物品的底层实现，
	// 这意味着返回值可以被断言为 DefaultItem
	UnderlyingItem() Item
	// NeedEnchOrRename 指示在导入这个
	// NBT 物品时是否需要附魔或重命名
	NeedEnchOrRename() bool
	// IsComplex 指示这个物品是否
	// 需要进一步的特殊处理才能得到
	IsComplex() bool
	// NeedCheckCompletely 指示在完成这个
	// NBT 物品的导入后是否需要检查其完整性。
	// 如果 NeedSpecialHandle 为假，
	// 则 NeedCheckCompletely 不应被使用
	NeedCheckCompletely() bool
	// NBTStableBytes 返回该物品在 NBT 部分的校验和。
	// NBT 的部分不包含物品名称和附魔数据，但包括物品
	// 组件和这个物品特定的一些 NBT 字段
	NBTStableBytes() []byte
	// TypeStableBytes 返回该种物品的种类哈希校验和。
	// 这意味着，同种的物品具有一致的种类哈希校验和
	TypeStableBytes() []byte
	// FullStableBytes 返回这个物品的哈希校验和
	FullStableBytes() []byte
}
