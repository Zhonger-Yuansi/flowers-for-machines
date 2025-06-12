package nbt_block

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	nbt_assigner_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/interface"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
	nbt_hash "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/hash"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/google/uuid"
)

func init() {
	nbt_assigner_interface.NBTBlockIsSupported = NBTBlockIsSupported
	nbt_assigner_interface.PlaceNBTBlock = PlaceNBTBlock
}

// NBTBlockIsSupported 检查 block 是否是受支持的 NBT 方块
func NBTBlockIsSupported(block nbt_parser_interface.Block) bool {
	switch block.(type) {
	case *nbt_parser_block.CommandBlock:
	case *nbt_parser_block.Sign:
	case *nbt_parser_block.StructureBlock:
	case *nbt_parser_block.Container:
	case *nbt_parser_block.Banner:
	case *nbt_parser_block.Frame:
	case *nbt_parser_block.Lectern:
	case *nbt_parser_block.JukeBox:
	default:
		return false
	}
	return true
}

// PlaceNBTBlock 根据传入的操作台和缓存命中系统，
// 在操作台的中心方块处制作一个 NBT 方块 nbtBlock。
//
// canFast 指示目标方块是否可以直接通过 setblock 放置。
//
// 如果不能通过 setblock 放置，那么 uniqueID 指示目标
// 方块所在结构的唯一标识，并且 offset 指示其相邻的可能
// 的方块，例如床的尾方块相对于头方块的偏移
func PlaceNBTBlock(
	console *nbt_console.Console,
	cache *nbt_cache.NBTCacheSystem,
	nbtBlock nbt_parser_interface.Block,
) (
	canFast bool,
	uniqueID uuid.UUID,
	offset protocol.BlockPos,
	err error,
) {
	// 初始化
	var method nbt_assigner_interface.Block
	hashNumber := nbt_hash.CompletelyHashNumber{
		HashNumber:    nbt_hash.NBTBlockHash(nbtBlock),
		SetHashNumber: nbt_hash.ContainerSetHash(nbtBlock),
	}

	// 检查是否可以快速放置
	if !nbtBlock.NeedSpecialHandle() {
		return true, uuid.UUID{}, protocol.BlockPos{}, nil
	}

	// 检查 NBT 缓存命中系统
	structure, hit, partHit := cache.NBTBlockCache().CheckCache(hashNumber)
	if hit && !partHit {
		return false, structure.UniqueID, structure.Offset, nil
	}

	// 取得相应 NBT 方块的放置方法
	switch block := nbtBlock.(type) {
	case *nbt_parser_block.CommandBlock:
		method = &CommandBlock{
			console: console,
			data:    *block,
		}
	case *nbt_parser_block.Sign:
		method = &Sign{
			console: console,
			data:    *block,
		}
	case *nbt_parser_block.StructureBlock:
		method = &StructrueBlock{
			console: console,
			data:    *block,
		}
	case *nbt_parser_block.Container:
		method = &Container{
			console: console,
			cache:   cache,
			data:    *block,
		}
	case *nbt_parser_block.Banner:
		method = &Banner{
			console: console,
			cache:   cache,
			data:    *block,
		}
	case *nbt_parser_block.Frame:
		method = &Frame{
			console: console,
			cache:   cache,
			data:    *block,
		}
	case *nbt_parser_block.Lectern:
		method = &Lectern{
			console: console,
			cache:   cache,
			data:    *block,
		}
	case *nbt_parser_block.JukeBox:
		method = &JukeBox{
			console: console,
			cache:   cache,
			data:    *block,
		}
	}

	// 放置相应方块
	err = method.Make()
	if err != nil {
		return false, uuid.UUID{}, protocol.BlockPos{}, fmt.Errorf("PlaceNBTBlock: %v", err)
	}

	// 检查完整性，如果需要的话
	if nbtBlock.NeedCheckCompletely() {
		nbtMap, err := simpleStructureGetter(console)
		if err != nil {
			return false, uuid.UUID{}, protocol.BlockPos{}, fmt.Errorf("PlaceNBTBlock: %v", err)
		}

		newBlock, err := nbt_parser_block.ParseNBTBlock(nbtBlock.BlockName(), nbtBlock.BlockStates(), nbtMap)
		if err != nil {
			return false, uuid.UUID{}, protocol.BlockPos{}, fmt.Errorf("PlaceNBTBlock: %v", err)
		}

		if hashNumber.HashNumber != nbt_hash.NBTBlockHash(newBlock) {
			return PlaceNBTBlock(console, cache, newBlock)
		}
	}

	// 保存缓存
	err = cache.NBTBlockCache().StoreCache(nbtBlock, method.Offset())
	if err != nil {
		return false, uuid.UUID{}, protocol.BlockPos{}, fmt.Errorf("PlaceNBTBlock: %v", err)
	}

	// 下次调用时将直接返回缓存
	return PlaceNBTBlock(console, cache, nbtBlock)
}
