package nbt_hash

import (
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
	nbt_parser_interface "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/interface"
	"github.com/cespare/xxhash/v2"
)

// NBTBlockHash 计算 block 的哈希校验和
func NBTBlockHash(block nbt_parser_interface.Block) uint64 {
	return xxhash.Sum64(block.StableBytes())
}

// NBTItemTypeHash 计算 item 的种类哈希校验和。
// 这意味着，对于两种相同的物品，它们具有相同的种类哈希校验和
func NBTItemTypeHash(item nbt_parser_interface.Item) uint64 {
	return xxhash.Sum64(item.TypeStableBytes())
}

// NBTItemFullHash 计算 item 的哈希校验和
func NBTItemFullHash(item nbt_parser_interface.Item) uint64 {
	return xxhash.Sum64(item.FullStableBytes())
}

// ContainerSetHash 计算 block 的集合哈希校验和。
// ContainerSetHash 假设给定的 block 可以断言为容器。
//
// 如果提供的 block 不是容器，或容器为空，
// 则返回 SetHashNumberNotExist (0)。
// 否则，返回这个容器的集合哈希校验和。
//
// 通常地，如果两个容器装有相同种类的物品，
// 且每个种类的物品数量相等，
// 则两个容器的集合哈希校验和相等
func ContainerSetHash(block nbt_parser_interface.Block) uint64 {
	container, ok := block.(*nbt_parser_block.Container)
	if !ok {
		return SetHashNumberNotExist
	}

	setBytes := container.SetBytes()
	if len(setBytes) == 0 {
		return SetHashNumberNotExist
	}

	return xxhash.Sum64(setBytes)
}
