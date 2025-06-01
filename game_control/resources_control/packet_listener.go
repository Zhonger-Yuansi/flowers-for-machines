package resources_control

import (
	"sync"

	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

// PacketListener 实现了一个可撤销监听的，
// 相对基础的数据包监听器
type PacketListener struct {
	mu           *sync.Mutex
	anyCallbacks []func(p packet.Packet) bool
	callbacks    map[uint32][]func(p packet.Packet) bool
}

// NewPacketListener 创建并返回一个新的 NewPacketListener
func NewPacketListener() *PacketListener {
	return &PacketListener{
		mu:           new(sync.Mutex),
		anyCallbacks: nil,
		callbacks:    make(map[uint32][]func(p packet.Packet) bool),
	}
}

// ListenPacket 监听数据包 ID 在 packetID 中的数据包，
// 并在收到这些数据包后执行回调函数 callback。
//
// 如果 callback 返回真，则此监听器将会被撤销，
// 否则将会继续保留。
//
// 如果 packetID 置空，则监听所有数据包
func (p *PacketListener) ListenPacket(
	packetID []uint32,
	callback func(p packet.Packet) bool,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(packetID) == 0 {
		p.anyCallbacks = append(p.anyCallbacks, callback)
		return
	}

	for _, pkID := range packetID {
		if p.callbacks[pkID] == nil {
			p.callbacks[pkID] = make([]func(p packet.Packet) bool, 0)
		}
		p.callbacks[pkID] = append(p.callbacks[pkID], callback)
	}
}

// onPacket ..
func (p *PacketListener) onPacket(pk packet.Packet) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Any callback
	{
		callbackCancel := make([]bool, len(p.anyCallbacks))
		haveAtLeastOneCancel := false
		for index, callbacks := range p.anyCallbacks {
			if callbacks(pk) {
				callbackCancel[index] = true
				haveAtLeastOneCancel = true
			}
		}

		if haveAtLeastOneCancel {
			newCallback := make([]func(p packet.Packet) bool, 0)
			for index, callbacks := range p.anyCallbacks {
				if callbackCancel[index] {
					continue
				}
				newCallback = append(newCallback, callbacks)
			}
			p.anyCallbacks = newCallback
		}
	}

	// Specific packet
	{
		packetID := pk.ID()
		cbs := p.callbacks[packetID]
		if cbs == nil {
			return
		}

		callbackCancel := make([]bool, len(cbs))
		haveAtLeastOneCancel := false
		for index, cb := range cbs {
			if cb(pk) {
				callbackCancel[index] = true
				haveAtLeastOneCancel = true
			}
		}

		if haveAtLeastOneCancel {
			newCallback := make([]func(p packet.Packet) bool, 0)

			for index, cb := range cbs {
				if callbackCancel[index] {
					continue
				}
				newCallback = append(newCallback, cb)
			}

			if len(newCallback) == 0 {
				delete(p.callbacks, packetID)
			} else {
				p.callbacks[packetID] = newCallback
			}
		}
	}
}
