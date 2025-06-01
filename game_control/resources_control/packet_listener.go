package resources_control

import (
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
	"github.com/google/uuid"
)

// singleListener 是单个数据包的监听器
type singleListener struct {
	uniqueID string                     // 该监听器的唯一标识符
	callback func(p packet.Packet) bool // 该监听器的回调函数
}

// PacketListener 实现了一个可撤销监听的，
// 相对基础的数据包监听器
type PacketListener struct {
	mu                      *sync.Mutex
	anyPacketListeners      []singleListener
	specificPacketListeners map[uint32][]singleListener
}

// NewPacketListener 创建并返回一个新的 NewPacketListener
func NewPacketListener() *PacketListener {
	return &PacketListener{
		mu:                      new(sync.Mutex),
		anyPacketListeners:      nil,
		specificPacketListeners: make(map[uint32][]singleListener),
	}
}

// ListenPacket 监听数据包 ID 在 packetID 中的数据包，
// 并在收到这些数据包后执行回调函数 callback。
//
// 如果 callback 返回真，则此监听器将会被撤销，
// 否则将会继续保留；
// 如果 packetID 置空，则监听所有数据包。
//
// 返回的 uniqueID 用于标识该监听器，以便于后续
// 调用 DestroyListener 以便于手动销毁监听器
func (p *PacketListener) ListenPacket(
	packetID []uint32,
	callback func(p packet.Packet) bool,
) (uniqueID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	uniqueID = uuid.New().String()
	listener := singleListener{
		uniqueID: uniqueID,
		callback: callback,
	}

	if len(packetID) == 0 {
		p.anyPacketListeners = append(p.anyPacketListeners, listener)
		return
	}

	for _, pkID := range packetID {
		if p.specificPacketListeners[pkID] == nil {
			p.specificPacketListeners[pkID] = make([]singleListener, 0)
		}
		p.specificPacketListeners[pkID] = append(p.specificPacketListeners[pkID], listener)
	}
	return
}

// DestroyListener 销毁唯一标识为 uniqueID 的数据包监听器。
// 如果这样的监听器不存在，则不会执行任何操作
func (p *PacketListener) DestroyListener(uniqueID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Any packet listener
	{
		found := false
		listenerIndex := 0

		for index, listener := range p.anyPacketListeners {
			if listener.uniqueID == uniqueID {
				found = true
				listenerIndex = index
				break
			}
		}

		if found {
			newListeners := make([]singleListener, 0)

			for index, listener := range p.anyPacketListeners {
				if index == listenerIndex {
					continue
				}
				newListeners = append(newListeners, listener)
			}

			p.anyPacketListeners = newListeners
			return
		}
	}

	// Specific packet listener
	for packetID, listeners := range p.specificPacketListeners {
		found := false
		listenerIndex := 0

		for index, listener := range listeners {
			if listener.uniqueID == uniqueID {
				found = true
				listenerIndex = index
				break
			}
		}

		if found {
			newListeners := make([]singleListener, 0)

			for index, listener := range listeners {
				if index == listenerIndex {
					continue
				}
				newListeners = append(newListeners, listener)
			}

			if len(newListeners) == 0 {
				delete(p.specificPacketListeners, packetID)
			} else {
				p.specificPacketListeners[packetID] = newListeners
			}

			return
		}
	}
}

// onPacket ..
func (p *PacketListener) onPacket(pk packet.Packet) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Any packet listener
	{
		callbackCancel := make([]bool, len(p.anyPacketListeners))
		haveAtLeastOneCancel := false
		for index, listeners := range p.anyPacketListeners {
			if listeners.callback(pk) {
				callbackCancel[index] = true
				haveAtLeastOneCancel = true
			}
		}

		if haveAtLeastOneCancel {
			newListeners := make([]singleListener, 0)
			for index, listeners := range p.anyPacketListeners {
				if callbackCancel[index] {
					continue
				}
				newListeners = append(newListeners, listeners)
			}
			p.anyPacketListeners = newListeners
		}
	}

	// Specific packet listener
	{
		packetID := pk.ID()
		listeners := p.specificPacketListeners[packetID]
		if listeners == nil {
			return
		}

		callbackCancel := make([]bool, len(listeners))
		haveAtLeastOneCancel := false
		for index, listener := range listeners {
			if listener.callback(pk) {
				callbackCancel[index] = true
				haveAtLeastOneCancel = true
			}
		}

		if haveAtLeastOneCancel {
			newListeners := make([]singleListener, 0)

			for index, listener := range listeners {
				if callbackCancel[index] {
					continue
				}
				newListeners = append(newListeners, listener)
			}

			if len(newListeners) == 0 {
				delete(p.specificPacketListeners, packetID)
			} else {
				p.specificPacketListeners[packetID] = newListeners
			}
		}
	}
}
