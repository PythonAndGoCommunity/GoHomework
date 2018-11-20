package json

import (
	"encoding/json"
)

func UnpackFromJSON(b []byte) *map[string]interface{}{
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil{
		panic(err)
	}
	return &m
}