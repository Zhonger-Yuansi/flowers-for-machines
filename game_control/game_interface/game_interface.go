package game_interface

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// ResourcesWrapper 是基于资源中心包装的机器人资源
type ResourcesWrapper struct {
	resources_control.BotInfo
	*resources_control.Resources
}

// GameInterface 实现了机器人与租赁服的高级交互，
// 例如基本的命令收发或高级的容器操作
type GameInterface struct {
	wrapper  *ResourcesWrapper
	commands *Commands
}

// NewResourcesWrapper 基于 resources 创建一个新的游戏交互器
func NewResourcesWrapper(resources *resources_control.Resources) *ResourcesWrapper {
	return &ResourcesWrapper{
		BotInfo:   resources.BotInfo(),
		Resources: resources,
	}
}

// NewGameInterface 基于 resources 创建一个新的游戏交互器
func NewGameInterface(resources *resources_control.Resources) *GameInterface {
	wrapper := NewResourcesWrapper(resources)
	return &GameInterface{
		wrapper:  wrapper,
		commands: NewCommands(wrapper),
	}
}

// GetBotInfo 返回机器人的基本信息
func (g *GameInterface) GetBotInfo() resources_control.BotInfo {
	return g.wrapper.BotInfo
}

// Commands 返回机器人在 MC 命令在收发上的相关实现
func (g *GameInterface) Commands() *Commands {
	return g.commands
}
