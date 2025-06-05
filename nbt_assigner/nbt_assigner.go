package nbt_assigner

import (
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/game_interface"
)

type NBTAssigner struct {
	mu  *sync.Mutex
	api *game_interface.GameInterface
}

// NewNBTAssigner 基于 api 创建并返回一个新的 NBTAssigner
func NewNBTAssigner(api *game_interface.GameInterface) *NBTAssigner {
	return &NBTAssigner{
		mu:  new(sync.Mutex),
		api: api,
	}
}
