package tcc

import "encoding/json"

// ObjToJson ...
func ObjToJson(v interface{}) string {
	value := ""
	switch s := v.(type) {
	case string:
		value = s
	default:
		b, err := json.Marshal(s)
		if err != nil {
			panic(err.Error())
		}
		value = string(b)
	}
	return value
}
