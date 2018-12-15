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
	QueryReg = regexp.MustCompile("^(get|GET|set|SET|del|DEL|keys|KEYS)")
	TopicReg = regexp.MustCompile("^(subscribe|SUBSCRIBE|publish|PUBLISH|unsubscribe|UNSUBSCRIBE)")
	ExitReg = regexp.MustCompile("^(exit|EXIT)$")
	DumpReg = regexp.MustCompile("^(dump|DUMP)$")
	RestoreReg = regexp.MustCompile("^(restore|RESTORE)$")
}
