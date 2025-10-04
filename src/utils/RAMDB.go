package utils

import "sync"

type RAMDB struct {
	data  map[string]any
	mutex sync.RWMutex
}

var ramDB RAMDB

func InitGlobalRAMDB() {
	ramDB = NewRAMDB()
}
func NewRAMDB() RAMDB {
	return RAMDB{
		data: make(map[string]any),
	}
}
func (t RAMDB) Set(key string, value any) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data[key] = value
}
func (t RAMDB) Get(key string) any {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data[key]
}
func (t RAMDB) Delete(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	delete(t.data, key)
}
