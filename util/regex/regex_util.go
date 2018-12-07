package regex

import (
	"regexp"
)

var (
	DoubleQuoteReg *regexp.Regexp
	QueryReg *regexp.Regexp
	TopicReg *regexp.Regexp
	ExitReg *regexp.Regexp
	DumpReg *regexp.Regexp
)

func init(){ 
	DoubleQuoteReg = regexp.MustCompile("\"(.*)\"")
	QueryReg = regexp.MustCompile("^(get|set|del|keys)")
	TopicReg = regexp.MustCompile("^(subscribe|publish|unsubscribe)")
	ExitReg = regexp.MustCompile("^exit$")
	DumpReg = regexp.MustCompile("^dump$")
}