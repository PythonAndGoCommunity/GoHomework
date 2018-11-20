package json

import (
	"encoding/json"
)

func PackToJSON(key string, value interface{}) []byte{
	m := map[string]interface{}{key : value}
	b, err := json.Marshal(m)
	if err != nil{
		panic(err)
	}
	return b
}

func PackMapToJson(m map[string]interface{}) []byte{
	b, err := json.Marshal(m)
	if err != nil{
		panic(err)
	}
	return b
}