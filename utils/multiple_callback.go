package utils

import "sync"

// MultipleCallback 包装在同一个事件上的多个回调函数
type MultipleCallback struct {
	mu        *sync.Mutex
	callbacks []func()
}

// NewMultipleCallback 创建一个新的 MultipleCallback
func NewMultipleCallback() *MultipleCallback {
	return &MultipleCallback{
		mu:        new(sync.Mutex),
		callbacks: nil,
	}
}

// Append 将一个新的回调函数加入底层切片
func (m *MultipleCallback) Append(f func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks = append(m.callbacks, f)
}

// FinishAll 执行底层切片的所有回调函数，并将切片清空
func (m *MultipleCallback) FinishAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, f := range m.callbacks {
		f()
	}
	m.callbacks = nil
}
