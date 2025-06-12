package nbt_assigner

import (
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_assigner_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/google/uuid"

	_ "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_block"
	_ "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_item"
	_ "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
	_ "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/item"
)

// NBTAssigner 是封装好的 NBT 方块放置实现
type NBTAssigner struct {
	mu      *sync.Mutex
	console *nbt_console.Console
	cache   *nbt_cache.NBTCacheSystem
}

// NewNBTAssigner 基于操作台和缓存命中系统创建并返回一个新的 NBTAssigner
func NewNBTAssigner(
	console *nbt_console.Console,
	cache *nbt_cache.NBTCacheSystem,
) *NBTAssigner {
	return &NBTAssigner{
		mu:      new(sync.Mutex),
		console: console,
		cache:   cache,
	}
}

// PlaceNBTBlock 制作一个新的 NBT 方块 nbtBlock。
// 制作位置是操作台的中心方块处。
//
// canFast 指示目标方块是否可以直接通过 setblock 放置。
//
// 如果不能通过 setblock 放置，那么 uniqueID 指示目标
// 方块所在结构的唯一标识，并且 offset 指示其相邻的可能
// 的方块，例如床的尾方块相对于头方块的偏移。
//
// PlaceNBTBlock 是阻塞的，它保证同一时刻只会制作一个
// NBT 方块
func (n *NBTAssigner) PlaceNBTBlock(nbtBlock nbt_parser_interface.Block) (
	canFast bool,
	uniqueID uuid.UUID,
	offset protocol.BlockPos,
	err error,
) {
	n.mu.Lock()
	defer n.mu.Unlock()
	canFast, uniqueID, offset, err = nbt_assigner_interface.PlaceNBTBlock(n.console, n.cache, nbtBlock)
	return
}
