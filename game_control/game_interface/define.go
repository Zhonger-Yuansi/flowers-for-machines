package game_interface

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/game_control/resources_control"
)

// 用于 PhoenixBuilder 与租赁服交互。
// 此结构体下的实现将允许您与租赁服进行交互操作，例如打开容器等
type GameInterface struct {
	botBasicInfo resources_control.BotBasicInfo
	resources    *resources_control.Resources
}
