package game_interface

import (
	"encoding/json"
	"fmt"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/google/uuid"
)

// Commands 是基于 ResourcesWrapper
// 实现的 MC 指令操作器，例如发送命令
// 并得到其响应体。
//
// 另外，出于对旧时代的尊重和可能的兼容性，
// 一些遗留实现也被同时迁移到此处
type Commands struct {
	api *ResourcesWrapper
}

// NewCommands 基于 api 创建并返回一个新的 Commands
func NewCommands(api *ResourcesWrapper) *Commands {
	return &Commands{api: api}
}

// packCommandRequest 根据给定的命令 command，
// 命令来源 origin 和命令请求 ID 包装一个命令请求体
func packCommandRequest(command string, origin uint32, requestID uuid.UUID) *packet.CommandRequest {
	return &packet.CommandRequest{
		CommandLine: command,
		CommandOrigin: protocol.CommandOrigin{
			Origin:    origin,
			UUID:      requestID,
			RequestID: "96045347-a6a3-4114-94c0-1bc4cc561694",
		},
		Internal:  false,
		UnLimited: false,
		Version:   0x24,
	}
}

// 向租赁服发送 Sizukana 命令且无视返回值。
// 当 dimensional 为真时，
// 将使用 execute 更换命令执行环境为机器人所在的环境
func (c *Commands) SendSettingsCommand(command string, dimensional bool) error {
	api := c.api

	if dimensional {
		command = fmt.Sprintf(
			`execute as @a[name="%s"] at @s run %s`,
			api.BotName,
			command,
		)
	}

	err := api.WritePacket(&packet.SettingsCommand{
		CommandLine:    command,
		SuppressOutput: true,
	})
	if err != nil {
		return fmt.Errorf("SendSettingsCommand: %v", err)
	}

	return nil
}

// sendCommand 以 origin 的身份向租赁服发送命令 command 并无视返回值
func (c *Commands) sendCommand(command string, origin uint32) error {
	err := c.api.WritePacket(
		packCommandRequest(
			command, origin, uuid.New(),
		),
	)
	if err != nil {
		return fmt.Errorf("sendCommand: %v", err)
	}
	return nil
}

// SendPlayerCommand 以玩家的身份向租赁服发送命令 command 并无视返回值
func (c *Commands) SendPlayerCommand(command string) error {
	err := c.sendCommand(command, protocol.CommandOriginPlayer)
	if err != nil {
		return fmt.Errorf("SendPlayerCommand: %v", err)
	}
	return nil
}

// SendPlayerCommand 以 Websocket 的身份向租赁服发送命令 command 并无视返回值
func (c *Commands) SendWSCommand(command string) error {
	err := c.sendCommand(command, protocol.CommandOriginAutomationPlayer)
	if err != nil {
		return fmt.Errorf("SendWSCommand: %v", err)
	}
	return nil
}

// sendCommandWithResp 以 origin 的身份向租赁服发送命令 command 并获取响应体
func (c *Commands) sendCommandWithResp(command string, origin uint32) (resp *packet.CommandOutput, err error) {
	api := c.api
	requestID := uuid.New()
	channel := make(chan struct{})

	api.Resources.Commands().SetCommandRequestCallback(
		requestID,
		func(p *packet.CommandOutput) {
			resp = p
			close(channel)
		},
	)

	err = api.WritePacket(
		packCommandRequest(
			command, origin, requestID,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("sendCommandWithResp: %v", err)
	}

	<-channel
	return
}

// SendPlayerCommandWithResp 以玩家的身份向租赁服发送命令 command 并获取响应体
func (c *Commands) SendPlayerCommandWithResp(command string) (resp *packet.CommandOutput, err error) {
	resp, err = c.sendCommandWithResp(command, protocol.CommandOriginPlayer)
	if err != nil {
		return nil, fmt.Errorf("SendPlayerCommandWithResp: %v", err)
	}
	return
}

// SendWSCommandWithResp 以 Websocket 的身份向租赁服发送命令 command 并获取响应体
func (c *Commands) SendWSCommandWithResp(command string) (resp *packet.CommandOutput, err error) {
	resp, err = c.sendCommandWithResp(command, protocol.CommandOriginAutomationPlayer)
	if err != nil {
		return nil, fmt.Errorf("SendWSCommandWithResp: %v", err)
	}
	return
}

// AwaitChangesGeneral 通过发送空指令以等待租赁服更改。
// 它曾被广泛使用而难以替代，但此处出于语义兼容性而保留
func (c *Commands) AwaitChangesGeneral() error {
	_, err := c.SendWSCommandWithResp("")
	if err != nil {
		return fmt.Errorf("AwaitChangesGeneral: %v", err)
	}
	return nil
}

// SendChat 使机器人在聊天栏说出 content 的内容
func (c *Commands) SendChat(content string) error {
	api := c.api

	err := api.WritePacket(
		&packet.Text{
			TextType:         packet.TextTypeChat,
			NeedsTranslation: false,
			SourceName:       api.BotName,
			Message:          content,
			XUID:             api.XUID,
			PlatformChatID:   "",
			Unknown1:         []string{"PlayerId", fmt.Sprintf("%d", api.EntityRuntimeID)},
		},
	)
	if err != nil {
		return fmt.Errorf("SendChat: %v", err)
	}

	return nil
}

// 以 actionbar 的形式向所有在线玩家显示 message
func (c *Commands) Title(message string) error {
	title := map[string]any{
		"rawtext": []any{
			map[string]any{
				"text": message,
			},
		},
	}
	jsonBytes, _ := json.Marshal(title)

	err := c.SendSettingsCommand(fmt.Sprintf("titleraw @a actionbar %s", jsonBytes), false)
	if err != nil {
		return fmt.Errorf("Title: %v", err)
	}

	return nil
}
