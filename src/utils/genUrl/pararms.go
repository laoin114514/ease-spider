package genUrl

import (
	"fmt"
	"sort"
)

func buildPararms(pararms map[string]any) string {
	if pararms == nil {
		return ""
	}
	keyArr := []string{}
	for k, _ := range pararms {
		keyArr = append(keyArr, k)
	}
	sort.Strings(keyArr)
	str := ""
	for _, key := range keyArr {
		str += fmt.Sprintf("%s=%v&", key, pararms[key])
	}
	str = str[:len(str)-1]
	return str
}
