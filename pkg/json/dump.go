package json

import "encoding/json"

func SafeDump(data interface{}) string {
	marshal, _ := json.MarshalIndent(data, "", "    ")
	return string(marshal)
}
