package handler

import (
	"strings"
	"NonRelDB/server/storage/inmemory"
	"NonRelDB/util/regex"
)

func HandleQuery(query string) string {
	queryParts := strings.Split(query, " ")[:2]

	switch queryCtx := strings.ToLower(queryParts[0]); queryCtx{
		case "get":{
			return inmemory.GetStorage().Get(queryParts[1])
		}
		case "set":{
			value := strings.Trim(regex.ValueReg.FindString(query),"\"")
			return inmemory.GetStorage().Set(queryParts[1], value)
		} 
		case "del":{
			return inmemory.GetStorage().Del(queryParts[1])
		}
		case "keys":{
			pattern := regex.ValueReg.FindString(query)
			return inmemory.GetStorage().Keys(pattern)
		}
	}

	return "Undefined query"
}
