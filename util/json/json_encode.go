package json

import (
	"encoding/json"
	"NonRelDB/log"
)

// PackToJSON receives string key and interface value and returns json bytes.
func PackToJSON(key string, value string) []byte{
	m := map[string]string{key : value}
	b, err := json.Marshal(m)
	if err != nil{
		log.Error.Println(err.Error())
	}
	return b
}

// PackToJSONIndent receives string key and interface value and returns json bytes with indent.
func PackToJSONIndent(key string, value string) []byte{
	m := map[string]string{key : value}
	b, err := json.MarshalIndent(m,"", " ")
	if err != nil{
		log.Error.Println(err.Error())
	}
	return b
}

// PackMapToJSON receives map and returns json bytes.
func PackMapToJSON(m map[string]string) []byte{
	b, err := json.Marshal(m)
	if err != nil{
		log.Error.Println(err.Error())
	}
	return b
}

// PackMapToJSONIndent receives map and returns json bytes with indent.
func PackMapToJSONIndent(m map[string]string) []byte{
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil{
		log.Error.Println(err.Error())
	}
	return b
}