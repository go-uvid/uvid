package tools

import "encoding/json"

func StructToJSONString(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
