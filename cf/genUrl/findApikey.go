package genUrl

import (
	"encoding/json"
	"os"
	"spider/typ"
)

func apiKey() []typ.ApiKey {
	var data []typ.ApiKey
	apiKey, _ := os.ReadFile("apiKey.json")
	json.Unmarshal(apiKey, &data)
	return data
}
func findApikey(handle string) (bool, typ.ApiKey) {
	apiKeys := apiKey()
	for _, apikey := range apiKeys {
		if apikey.Name == handle {
			return true, apikey
		}
	}
	return false, typ.ApiKey{}
}
