package block_helper

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
)

// ContainerBlockHelper 描述了一个容器，
// 并记载了它的坐标、名称和方块状态。
//
// ConsiderOpenDirection 指示打开目标容器
// 是否需要考虑其打开方向上的障碍物方块，这
// 似乎只对箱子和潜影盒有效。
// 当其为真时，应在 Facing 填写该容器的朝向，
// 否则可以置为默认的零值。
//
// Content 指示这个容器装有的物品，
// 为了方便，此处总是使用定长度的数组，
// 但尝试修改超出目标容器格子数的格子
// 会导致程序的行为未定义
type ContainerBlockHelper struct {
	Name   string
	States map[string]any

	ConsiderOpenDirection bool
	Facing                uint8

	Content [27]*protocol.ItemInstance
}

func (c ContainerBlockHelper) BlockName() string {
	return c.Name
}

func (c ContainerBlockHelper) BlockStates() map[string]any {
	return c.States
}

func (c ContainerBlockHelper) BlockStatesString() string {
	return utils.MarshalBlockStates(c.States)
}

// ShouldCleanNearBlock 指示打开该容器前是否需要清除
// 其相邻的方块。offset 指示这个相邻方块的位置。这目前
// 只对箱子和潜影盒有用
func (c ContainerBlockHelper) ShouldCleanNearBlock() (offset [3]int32, needClean bool) {
	if !c.ConsiderOpenDirection {
		return [3]int32{}, false
	}

	switch c.Facing {
	case 0:
		return [3]int32{0, -1, 0}, true
	case 1:
		return [3]int32{0, 1, 0}, true
	case 2:
		return [3]int32{0, 0, -1}, true
	case 3:
		return [3]int32{0, 0, 1}, true
	case 4:
		return [3]int32{-1, 0, 0}, true
	case 5:
		return [3]int32{1, 0, 0}, true
	}

	return
}

// Contents 返回该容器装有的物品。
// 如果一个容器的格子数量低于 27 格，
// 那么尝试修改超出格子数量部分的物品
// 会导致未定义行为
func (c *ContainerBlockHelper) Contents() [27]*protocol.ItemInstance {
	return c.Content
}
