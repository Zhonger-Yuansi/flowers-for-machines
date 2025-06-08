package item_cache

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/block_helper"
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
	"github.com/google/uuid"
)

// ------------------------- Item information -------------------------

// ItemCacheInfo 记载单个物品缓存的路径、
// 数量以及该物品的完整校验和与集合校验和
type ItemCacheInfo struct {
	// 该物品所在的槽位索引
	SlotID resources_control.SlotID
	// 该物品的数量
	Count uint8
	// Hash 是这个物品的哈希校验和
	Hash nbt_hash.CompletelyHashNumber
}

// CompletelyItemInfo 在 ItemCacheInfo 基础上描述了这个物品所在容器的信息
type CompletelyItemInfo struct {
	// 我们应当如何打开这个容器
	ContainerInfo block_helper.ContainerBlockOpenInfo
	// 该物品在的校验和与位置信息
	ItemInfo ItemCacheInfo
}

// ------------------------- Structure -------------------------

// StructureItemCache 指示一个被储存在结构中的缓存物品，
// 且结构保存的是一个容器方块，而缓存物品便在其中
type StructureItemCache struct {
	// UniqueID 是这个结构的唯一标识符
	UniqueID uuid.UUID
	// CompletelyInfo 是这个缓存物品的信息
	CompletelyInfo CompletelyItemInfo
}

// StructureItems 记载了一个结构中的容器方块，
// 而这个容器方块包含了一些物品。
// 除此外，StructureItems 还记载该容器中每个物
// 品的完备信息，以便于实际操作时可以保持操作台
// 中容器帮助方块中的一致性
type StructureItems struct {
	// 我们应当如何打开这个容器
	ContainerInfo block_helper.ContainerBlockOpenInfo
	// 该容器所包含的物品
	Items []CompletelyItemInfo
}

// ------------------------- End -------------------------
