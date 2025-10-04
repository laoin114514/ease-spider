package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type JSONDB struct {
	data  map[string]any
	mutex sync.RWMutex
	path  string
}

var JsonDB *JSONDB

func InitGlobalJSONDB(path string) {
	JsonDB = NewJSONDB(path)
}
func NewJSONDB(path string) *JSONDB {
	return &JSONDB{
		data: make(map[string]any),
		path: path,
	}
}

func (j *JSONDB) getOldData() map[string]any {
	oldDataStr, err := os.ReadFile(j.path)
	if err != nil {
		j.writeFile(map[string]any{})
	}
	var oldData = map[string]any{}
	json.Unmarshal(oldDataStr, &oldData)
	return oldData
}
func (j *JSONDB) writeFile(oldData map[string]any) {
	file, err := os.OpenFile(j.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(oldData)
}
func (j *JSONDB) Set(key string, value any) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.data = j.getOldData()
	j.data[key] = value
	j.writeFile(j.data)

}
func (j *JSONDB) Get(key string) any {
	j.mutex.RLock()
	defer j.mutex.RUnlock()
	j.data = j.getOldData()
	return j.data[key]
}
func (j *JSONDB) Delete(key string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.data = j.getOldData()
	delete(j.data, key)
}
