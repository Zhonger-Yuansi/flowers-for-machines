package item_stack_operation

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
)

type Default struct{}

func (Default) ID() uint8 {
	return IDItemStackOperationDefault
}

func (Default) CanInline() bool {
	return false
}

func (Default) Make(runtiemData MakingRuntime) []protocol.StackRequestAction {
	return nil
}
