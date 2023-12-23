package mapwmutex

import (
	"sync"
)

type MapWmutex[Hash comparable, Value any] struct {
	data map[Hash]Value
	m    *sync.RWMutex
}

func (mWm *MapWmutex[Hash, Value]) Store(key Hash, val Value) {
	mWm.m.Lock()
	mWm.data[key] = val
	mWm.m.Unlock()
}

func (mWm *MapWmutex[Hash, Value]) Load(key Hash) Value {
	mWm.m.RLock()
	val := mWm.data[key]
	mWm.m.RUnlock()

	return val
}

func (mWm *MapWmutex[Hash, Value]) LoadWithStatus(key Hash) (val Value, ok bool) {
	mWm.m.RLock()
	val, ok = mWm.data[key]
	mWm.m.RUnlock()

	return
}

func (mWm *MapWmutex[Hash, Value]) RowCount() (length int) {
	mWm.m.RLock()
	length = len(mWm.data)
	mWm.m.RUnlock()

	return length
}

func (mWm *MapWmutex[Hash, Value]) Delete(key Hash) {
	mWm.m.Lock()
	delete(mWm.data, key)
	mWm.m.Unlock()
}

func (mWm *MapWmutex[Hash, Value]) Reset() {
	mWm.m.Lock()
	mWm.data = make(map[Hash]Value)
	mWm.m.Unlock()
}

func (mWm *MapWmutex[Hash, Value]) GetData() map[Hash]Value {
	mWm.m.RLock()

	copyMap := make(map[Hash]Value, len(mWm.data))
	for key, value := range mWm.data {
		copyMap[key] = value
	}

	mWm.m.RUnlock()

	return copyMap
}

func NewMapWmutex[Hash comparable, Value any](len int) *MapWmutex[Hash, Value] {
	return &MapWmutex[Hash, Value]{
		data: make(map[Hash]Value, len),
		m:    new(sync.RWMutex),
	}
}
