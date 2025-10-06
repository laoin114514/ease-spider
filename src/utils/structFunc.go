package utils

import (
	"encoding/json"
	"fmt"
	"sort"
)

type StructFunc struct {
	v any
}

func NewStructFunc(v any) *StructFunc {
	return &StructFunc{
		v: v,
	}
}

// 将结构体转换为map
func (s *StructFunc) StructToMap() (map[string]any, error) {
	data, err := json.Marshal(s.v)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 将结构体转换为有序参数字符串
func (s *StructFunc) StructToOrderParams() (string, error) {
	mapData, err := s.StructToMap()
	if err != nil {
		return "", err
	}
	keyArr := []string{}
	for k, _ := range mapData {
		keyArr = append(keyArr, k)
	}
	sort.Strings(keyArr)
	str := ""
	for _, key := range keyArr {
		str += fmt.Sprintf("%s=%v&", key, mapData[key])
	}
	str = str[:len(str)-1]
	return str, nil
}
