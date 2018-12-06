package regex

import (
	"regexp"
)

var (
	DoubleQuoteReg *regexp.Regexp
	QueryReg *regexp.Regexp
	ExitReg *regexp.Regexp
	DumpReg *regexp.Regexp
)

func init(){ 
	DoubleQuoteReg = regexp.MustCompile("\"(.*)\"")
	QueryReg = regexp.MustCompile("^(get|set|del|keys)")
	ExitReg = regexp.MustCompile("^exit$")
	DumpReg = regexp.MustCompile("^dump$")
}