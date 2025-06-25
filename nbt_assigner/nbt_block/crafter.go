package nbt_block

import (
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_cache"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_assigner/nbt_console"
	nbt_parser_block "github.com/Happy2018new/the-last-problem-of-the-humankind/nbt_parser/block"
)

// 合成器
type Crafter struct {
	console *nbt_console.Console
	cache   *nbt_cache.NBTCacheSystem
	data    nbt_parser_block.Crafter
}

func (Crafter) Offset() protocol.BlockPos {
	return protocol.BlockPos{0, 0, 0}
}

func (c *Crafter) AsContainer() *Container {
	return &Container{
		console: c.console,
		cache:   c.cache,
		data:    *c.data.AsContainer(),
	}
}

func (c *Crafter) Make() error {
	api := c.console.API()
	center := c.console.Center()

	// 处理容器内的物品
	err := c.AsContainer().Make()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	// 如果这个合成台没有被禁用的物品栏，
	// 则可以直接返回值
	if c.data.NBT.DisabledSlots == 0 {
		return nil
	}

	// 否则，开始设置被禁用的物品栏
	for index := range 9 {
		if c.data.NBT.DisabledSlots&int16(1<<index) == 0 {
			continue
		}
		err = api.Resources().WritePacket(&packet.PlayerToggleCrafterSlotRequest{
			PosX:     center[0],
			PosY:     center[1],
			PosZ:     center[2],
			Slot:     byte(index),
			Disabled: true,
		})
		if err != nil {
			return fmt.Errorf("Make: %v", err)
		}
	}

	// 等待更改
	err = api.Commands().AwaitChangesGeneral()
	if err != nil {
		return fmt.Errorf("Make: %v", err)
	}

	return nil
}
