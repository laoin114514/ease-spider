package genUrl

import (
	"encoding/json"
	"os"
)

type ApiKey struct {
	Name   string `json:"name"`
	ApiKey string `json:"apikey"`
	Secret string `json:"secret"`
}

func apiKey() []ApiKey {
	var data []ApiKey
	apiKey, _ := os.ReadFile("apiKey.json")
	json.Unmarshal(apiKey, &data)
	return data
}
func findApikey(handle string) (bool, ApiKey) {
	apiKeys := apiKey()
	for _, apikey := range apiKeys {
		if apikey.Name == handle {
			return true, apikey
		}
	}
	return false, ApiKey{}
}
