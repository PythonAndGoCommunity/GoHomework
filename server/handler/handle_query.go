package handler

import (
	"NonRelDB/util/sync"
	"strings"
	"NonRelDB/util/regex"
)

// HandleQuery handling queries to db.
func HandleQuery(query string, syncMap *sync.Map) string {
	queryParts := strings.Split(query, " ")
	
	if len (queryParts) >= 2 {
		switch queryCtx := strings.ToLower(queryParts[0]); queryCtx{
			case "get":{
				return syncMap.Get(queryParts[1])
			}
			case "set":{
				value := strings.Trim(regex.DoubleQuoteReg.FindString(query),"\"")
				return syncMap.Set(queryParts[1], value)
			} 
			case "del":{
				return syncMap.Del(queryParts[1])
			}
			case "keys":{
				pattern := strings.Trim(regex.DoubleQuoteReg.FindString(query),"\"")
				return syncMap.Keys(pattern)
			}
		}
	}

	return "Undefined query"
}
