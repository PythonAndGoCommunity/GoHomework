package regex

import (
	"regexp"
)

var ValueReg *regexp.Regexp
var QueryReg *regexp.Regexp
var ExitReg *regexp.Regexp
var DumpReg *regexp.Regexp

func init(){ 
	ValueReg = regexp.MustCompile("\"(.*)\"")
	QueryReg = regexp.MustCompile("^(get|set|del|keys)")
	ExitReg = regexp.MustCompile("^exit$")
	DumpReg = regexp.MustCompile("^dump$")
}