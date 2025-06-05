package block_helper

// BlockHelper 是工作台上的帮助类方块，
// 但也可以表述操控台中心的方块或中心方
// 块及帮助类方块的相邻方块。
//
// 如果表示的是一个帮助类方块，
// 那么它可以是容器、铁砧或织布机
type BlockHelper interface {
	// BlockName 获取该方块的名称
	BlockName() string
	// BlockStates 获取该方块的方块状态
	BlockStates() map[string]any
	// BlockStatesString 获取该方块的方块状态的字符串表示
	BlockStatesString() string
}
