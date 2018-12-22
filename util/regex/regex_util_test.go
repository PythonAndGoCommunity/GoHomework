package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleQuoteReg__Matching__Success(t *testing.T) {
	assert.True(t, DoubleQuoteReg.MatchString("\"123\""))
}

func TestDoublQuoteReg__NotMatching__Fail(t *testing.T) {
	assert.False(t, DoubleQuoteReg.MatchString("123"))
}

func TestDoubleQuoteReg__Parse123__Success(t *testing.T) {
	assert.Equal(t, "\"123\"", DoubleQuoteReg.FindString("\"123\""))
}

func TestDoubleQuoteReg__Parse123WithoutQuotes__Fail(t *testing.T) {
	assert.NotEqual(t, "\"123\"", DoubleQuoteReg.FindString("123"))
}

func TestQueryReg__Get__Success(t *testing.T) {
	assert.True(t, QueryReg.MatchString("get 123"))
}

func TestQueryReg__Set__Success(t *testing.T) {
	assert.True(t, QueryReg.MatchString("get 123 \"123\""))
}

func TestQueryReg__Del__Success(t *testing.T) {
	assert.True(t, QueryReg.MatchString("get 123"))
}

func TestQueryReg__Keys__Success(t *testing.T) {
	assert.True(t, QueryReg.MatchString("keys \"\\*\""))
}

func TestQueryReg__123__Fail(t *testing.T) {
	assert.False(t, QueryReg.MatchString("123"))
}

func TestTopicReg__Subscribe__Success(t *testing.T) {
	assert.True(t, TopicReg.MatchString("subscribe redis"))
}

func TestTopicReg__Unsubscribe__Success(t *testing.T) {
	assert.True(t, TopicReg.MatchString("unsubscribe redis"))
}

func TestTopicReg__Publish__Success(t *testing.T) {
	assert.True(t, TopicReg.MatchString("publish redis \"Hello World\""))
}

func TestTopicReg__123__Fail(t *testing.T) {
	assert.False(t, TopicReg.MatchString("123"))
}

func TestExitReg__Exit__Success(t *testing.T) {
	assert.True(t, ExitReg.MatchString("exit"))
}

func TestExitReg__ExitWithSpaces__Fail(t *testing.T) {
	assert.False(t, ExitReg.MatchString(" exit "))
}

func TestDumpReg__Dump__Success(t *testing.T) {
	assert.True(t, DumpReg.MatchString("dump"))
}

func TestDumpReg__DumpWithSpaces__Fail(t *testing.T) {
	assert.False(t, DumpReg.MatchString(" dump "))
}

func TestRestoreReg__Restore__Success(t *testing.T) {
	assert.True(t, RestoreReg.MatchString("restore"))
}

func TestRestoreReg__RestoreWithSpaces__Fail(t *testing.T) {
	assert.False(t, RestoreReg.MatchString(" restore "))
}
