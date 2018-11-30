package handler

import (
	"strings"
	"NonRelDB/server/storage/inmemory"
	"NonRelDB/util/json"
)

func HandleQuery(q string) string {
	p := strings.Split(q, " ")
	if len(p) == 1 {
		if strings.ToLower(p[0]) == "dump"{
			return string(json.PackMapToJSON((*inmemory.GetStorage().GetMap())))
		}
	}
	if len(p) == 2 {
		if strings.ToLower(p[0]) == "get" {
			return inmemory.GetStorage().Get(p[1])
		} else if strings.ToLower(p[0]) == "del" {
			return inmemory.GetStorage().Del(p[1])
		} else if strings.ToLower(p[0]) == "keys" {
			return inmemory.GetStorage().Keys(p[1])
		}
	} else if len(p) >= 3 {
		if strings.ToLower(p[0]) == "set" {
			return inmemory.GetStorage().Set(p[1], strings.Join(p[2:]," "))
		}
	}
	return "Undefined query"
}
