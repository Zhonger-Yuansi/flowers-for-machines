package resources_control

import (
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

const (
	ContainerStatesHaveNotOpen uint8 = iota
	ContainerStatesOpening
	ContainerStatesClosed
)

// ContainerManager 描述一个在内存中维护的容器实现，
// 它用于追踪和监控目前已打开容器的状态
type ContainerManager struct {
	mu     *sync.Mutex
	states uint8

	openingData  *packet.ContainerOpen
	openCallback func()

	closingData   *packet.ContainerClose
	closeCallback func(isServerSide bool)
}

// NewContainerManager 创建并返回一个新的容器管理器
func NewContainerManager() *ContainerManager {
	return &ContainerManager{
		mu:            new(sync.Mutex),
		openingData:   nil,
		openCallback:  nil,
		closingData:   nil,
		closeCallback: nil,
	}
}

// States 返回已打开容器的状态。
// 目前只存在 3 种状态：
//   - 0: 曾经没有打开过容器
//   - 1: 目前存在一个已被打开的容器
//   - 2: 曾经打开过容器，但是关闭了
func (c ContainerManager) States() uint8 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.states
}

// ContainerData 获取当前已打开容器的状态。
// 返回的 existed 指示当前是否已经打开了容器
func (c ContainerManager) ContainerData() (data packet.ContainerOpen, existed bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.states != ContainerStatesOpening {
		return packet.ContainerOpen{}, false
	}
	return *c.openingData, true
}

// SetContainerOpenCallback 设置容器打开时应该执行的回调函数。
// 另外，设置的回调函数会在其被执行后被移除
func (c *ContainerManager) SetContainerOpenCallback(f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.openCallback = f
}

// SetContainerCloseCallback 设置容器关闭时应该执行的回调函数。
// isServerSide 指示容器是否是由服务器强制关闭的。
// 另外，设置的回调函数会在其被执行后被移除
func (c *ContainerManager) SetContainerCloseCallback(f func(isServerSide bool)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeCallback = f
}

// onContainerOpen ..
func (c *ContainerManager) onContainerOpen(p *packet.ContainerOpen) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.openingData = p
	c.closingData = nil
	c.states = ContainerStatesOpening

	if c.openCallback != nil {
		go c.openCallback()
		c.openCallback = nil
	}
}

// ContainerClose ..
func (c *ContainerManager) onContainerClose(p *packet.ContainerClose) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closingData = p
	c.openingData = nil
	c.states = ContainerStatesClosed

	if c.closeCallback != nil {
		go c.closeCallback(p.ServerSide)
		c.closeCallback = nil
	}
}
