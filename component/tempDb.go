package component

import (
	"encoding/json"
	"fmt"
	"os"
)

var tempDB TempDB

type TempDB struct {
}

func getOldData() map[string]any {
	oldDataStr, err := os.ReadFile("tempDB.json")
	if err != nil {
		writeFile(map[string]any{})
	}
	var oldData = map[string]any{}
	json.Unmarshal(oldDataStr, &oldData)
	return oldData
}
func writeFile(oldData map[string]any) {
	file, err := os.OpenFile("tempDB.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(oldData)
}
func (t TempDB) Set(key string, value any) {
	oldData := getOldData()
	oldData[key] = value
	writeFile(oldData)
}
func (t TempDB) Get(key string) any {
	data := getOldData()
	return data[key]
}
func (t TempDB) Delete(key string) {
	oldData := getOldData()
	delete(oldData, key)
	writeFile(oldData)
}
