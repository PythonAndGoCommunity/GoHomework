package handler

import (
	"fmt"
	"strings"
	"regexp"
	"NonRelDB/server/storage/inmemory"
)

var keyGetDelReg *regexp.Regexp
var keySetReg *regexp.Regexp
var valueReg *regexp.Regexp
var getReg *regexp.Regexp
var setReg *regexp.Regexp
var delReg *regexp.Regexp
var keysReg *regexp.Regexp

func init(){
	
	keyGetDelReg = regexp.MustCompile("\\s(.*)$")

	keySetReg = regexp.MustCompile("\\s(.*)\\s")

	valueReg = regexp.MustCompile("\"(.*)\"$")
	
	getReg = regexp.MustCompile("^get\\s(.*)$")

	setReg = regexp.MustCompile("^set\\s(.*)\\s\"(.*)\"$")

	delReg = regexp.MustCompile("^del\\s(.*)$")

	keysReg = regexp.MustCompile("^keys\\s\"(.*?)\"$")
}

func HandleQuery(q string) string {
	fmt.Println("[" + q + "]")

	if getReg.MatchString(q) {
		key := strings.Trim(keyGetDelReg.FindString(q), " ")
		fmt.Println("[" + key + "]")
		return inmemory.GetStorage().Get(key)

	} else if setReg.MatchString(q) {
		key := strings.Trim(keySetReg.FindString(q), " ")
		value := strings.Trim(valueReg.FindString(q), "\"")
		fmt.Println("[" + key + "]")
		fmt.Println("[" + value + "]")
		return inmemory.GetStorage().Set(key, value )

	} else if delReg.MatchString(q) {
		key := strings.Trim(keyGetDelReg.FindString(q), " ")
		fmt.Println("[" + key + "]")
		return inmemory.GetStorage().Del(key)
		
	} else if keysReg.MatchString(q) {
		reg := strings.Trim(valueReg.FindString(q),"\"")
		fmt.Println("[" + reg + "]")
		return inmemory.GetStorage().Keys(reg)

	}

	return "Undefined query"
}
