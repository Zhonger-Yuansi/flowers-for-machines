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

func NBTBlockIsSupported(block nbt_parser_interface.Block) bool {
	switch block.(type) {
	case *nbt_parser_block.CommandBlock:
	case *nbt_parser_block.Sign:
	case *nbt_parser_block.StructureBlock:
	case *nbt_parser_block.Container:
	default:
		return false
	}
	return true
}

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
	nbtBlockBase := NBTBlockBase{
		console: console,
		cache:   cache,
	}
	hashNumber := nbt_hash.CompletelyHashNumber{
		HashNumber:    nbt_hash.NBTBlockHash(nbtBlock),
		SetHashNumber: nbt_hash.ContainerSetHash(nbtBlock),
	}

	// 检查是否可以快速放置
	if !nbtBlock.NeedSpecialHandle() {
		return true, uuid.UUID{}, protocol.BlockPos{}, nil
	}

	// 检查 NBT 缓存命中系统
	structure, hit, _ := cache.NBTBlockCache().CheckCache(hashNumber)
	if hit {
		return false, structure.UniqueID, structure.Offset, nil
	}

	// 取得相应 NBT 方块的放置方法
	switch block := nbtBlock.(type) {
	case *nbt_parser_block.CommandBlock:
		method = &CommandBlock{
			NBTBlockBase: nbtBlockBase,
			data:         *block,
		}
	case *nbt_parser_block.Sign:
		method = &Sign{
			NBTBlockBase: nbtBlockBase,
			data:         *block,
		}
	case *nbt_parser_block.StructureBlock:
		method = &StructrueBlock{
			NBTBlockBase: nbtBlockBase,
			data:         *block,
		}
	case *nbt_parser_block.Container:
		method = &Container{
			NBTBlockBase: nbtBlockBase,
			data:         *block,
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
