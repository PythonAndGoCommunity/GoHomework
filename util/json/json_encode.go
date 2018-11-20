package json

import (
	"encoding/json"
)

// PackToJSON | Receives string key and interface value and returns json bytes.
func PackToJSON(key string, value interface{}) []byte{
	m := map[string]interface{}{key : value}
	b, err := json.MarshalIndent(m,""," ")
	if err != nil{
		panic(err)
	}
	return b
}

// PackMapToJSON | Receives map and returns json bytes.
func PackMapToJSON(m map[string]interface{}) []byte{
	b, err := json.MarshalIndent(m,""," ")
	if err != nil{
		panic(err)
	}
	return b
}