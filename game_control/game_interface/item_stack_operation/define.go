package item_stack_operation

import "github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"

const (
	IDItemStackOperationDefault uint8 = iota
	IDItemStackOperationMove
	IDItemStackOperationSwap
	IDItemStackOperationDrop
	IDItemStackOperationCreativeItem
	IDItemStackOperationHighLevelRenaming
	IDItemStackOperationHighLevelLooming
)

// ItemStackOperation 指示所有实现了它的物品操作
type ItemStackOperation interface {
	// CanInline 指示该物品操作是否可以内联到单个物品堆栈操作请求中。
	// 如果不能，则该物品操作则应被内联到同一个数据包的另外一个请求中。
	//
	// 在部分特殊情况下，如内联操作与非内联操作相邻，则非内联操作应当
	// 被接下来相邻的所有非内联操作放置在一个新的数据包中；
	// 亦或，接下来相邻的所有内联操作被放置在一个新的数据包中。
	//
	// 最终，构造的数据包应当满足：
	// {
	//		...
	//		<[can inline, can inline, ..]>,
	//		<[can't inline], [can't inline], ...>,
	//		<[can inline, can inline, ..]>,
	//		<[can't inline], [can't inline], ...>,
	//		...
	// }
	//
	// 其中，每个尖括号表示一个 packet.ItemStackRequest 数据包，
	// 每个方括号表示一个单个的物品堆栈请求
	CanInline() bool
	// ID 指示该物品操作的编号，它是自定义的
	ID() uint8
	// Make 基于运行时数据 runtiemData，
	// 返回目标物品操作的标准物品堆栈请求的动作
	Make(runtiemData MakingRuntime) []protocol.StackRequestAction
}
