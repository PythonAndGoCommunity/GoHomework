package json

import (
	"encoding/json"
	"NonRelDB/server/log"
)

// PackToJSON | Receives string key and interface value and returns json bytes.
func PackToJSON(key string, value string) []byte{
	m := map[string]string{key : value}
	b, err := json.MarshalIndent(m,""," ")
	if err != nil{
		log.Error.Println(err.Error())
	}
	return b
}

// PackMapToJSON | Receives map and returns json bytes.
func PackMapToJSON(m map[string]string) []byte{
	b, err := json.MarshalIndent(m,""," ")
	if err != nil{
		log.Error.Println(err.Error())
	}
	return b
}