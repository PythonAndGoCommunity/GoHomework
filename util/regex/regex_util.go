package regex

import (
	"regexp"
)

var (
	// DoubleQuoteReg regexp for values inside double quotes.
	DoubleQuoteReg *regexp.Regexp
	// QueryReg regexp which checks is this query to db.
	QueryReg *regexp.Regexp
	// TopicReg regexp which checks is this topic's command.
	TopicReg *regexp.Regexp
	// ExitReg regexp which checks is this exit command.
	ExitReg *regexp.Regexp
	// DumpReg regexp which checks is this dump command.
	DumpReg *regexp.Regexp
	// RestoreReg regexp which check is this restore command.
	RestoreReg *regexp.Regexp
)

func init() {
	DoubleQuoteReg = regexp.MustCompile("\"(.*)\"")
	QueryReg = regexp.MustCompile("^(get|set|del|keys)")
	TopicReg = regexp.MustCompile("^(subscribe|publish|unsubscribe)")
	ExitReg = regexp.MustCompile("^exit$")
	DumpReg = regexp.MustCompile("^dump$")
	RestoreReg = regexp.MustCompile("^restore$")
}
