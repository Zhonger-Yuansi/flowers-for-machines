package nbt_block

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
)

type NBTBlockBase struct {
	console *nbt_console.Console
	cache   *nbt_cache.NBTCacheSystem
}
