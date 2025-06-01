package resources_control

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/utils"
	"github.com/google/uuid"
)

// CommandRequestCallback 是简单的指令回调维护器
type CommandRequestCallback struct {
	callback utils.SyncMap[uuid.UUID, func(p *packet.CommandOutput)]
}

// NewCommandRequestCallback 创建并返回一个新的 CommandRequestCallback
func NewCommandRequestCallback() *CommandRequestCallback {
	return new(CommandRequestCallback)
}

// SetCommandRequestCallback 设置当收到请求 ID 为 requestID 的命令请求的响应后，
// 应当执行的回调函数 f。其中，p 指示服务器发送的针对此命令请求的响应体
func (c *CommandRequestCallback) SetCommandRequestCallback(
	requestID uuid.UUID,
	f func(p *packet.CommandOutput),
) {
	c.callback.Store(requestID, f)
}

// DeleteCommandRequestCallback 清除请求
// ID 为 requestID 的命令请求的回调函数。
// 此函数应当只在命令请求超时的时候被调用
func (c *CommandRequestCallback) DeleteCommandRequestCallback(requestID uuid.UUID) {
	c.callback.Delete(requestID)
}

// onCommandOutput ..
func (c *CommandRequestCallback) onCommandOutput(p *packet.CommandOutput) {
	cb, ok := c.callback.LoadAndDelete(p.CommandOrigin.UUID)
	if ok {
		cb(p)
	}
}
