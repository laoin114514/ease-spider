package newLuogu

import (
	"encoding/json"
	"log"
	"math"
	"net/url"
)

func decode(jsonEncode string) int {
	decoded, err := url.QueryUnescape(jsonEncode)
	if err != nil {
		log.Fatal("解码错误:", err)
	}

	// 解析JSON到结构体
	var data map[string]any
	err = json.Unmarshal([]byte(decoded), &data)
	if err != nil {
		log.Fatal("JSON解析错误:", err)
	}
	if data["code"].(float64) != 200 {
		return 0
	}
	count := data["currentData"].(map[string]any)["records"].(map[string]any)["count"].(float64)
	perPage := data["currentData"].(map[string]any)["records"].(map[string]any)["perPage"].(float64)
	page := int(math.Ceil(count / perPage))
	return page
}
func countPage(data map[string]any) int {
	count := data["count"].(float64)
	perPage := data["perPage"].(float64)
	page := int(math.Ceil(count / perPage))
	return page
}
