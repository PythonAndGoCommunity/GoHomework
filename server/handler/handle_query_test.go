package handler

import (
	"NonRelDB/util/sync"
	"testing"
	"github.com/stretchr/testify/assert"
)


var testSyncMap *sync.Map

// Used as setup in testing.
func init(){
	testMap := map[string]string{"1" : "one", "2" : "two", "3" : "three"}
	testSyncMap = &sync.Map{}
	testSyncMap.SetMap(&testMap)
}

func TestHandleQuery__GetExisting__Found(t *testing.T){
	assert.Equal(t, "one", HandleQuery("get 1", testSyncMap))
}

func TestHandleQuery__GetExisting__NotFound(t *testing.T){
	assert.Equal(t, "Value with this key not found", HandleQuery("get 4", testSyncMap))
}

func TestHandleQuery__SetExisting__Changed(t *testing.T){
	assert.Equal(t, "two", HandleQuery("get 2", testSyncMap))
	assert.Equal(t, "Value has changed", HandleQuery("set 2 \"2\"",testSyncMap))
	assert.Equal(t, "2", HandleQuery("get 2", testSyncMap))
}

func TestHandleQuery__SetNonExistring__Created(t *testing.T){
	assert.Equal(t, "Value with this key not found", HandleQuery("get 4", testSyncMap))
	assert.Equal(t, "Value has changed", HandleQuery("set 4 \"four\"", testSyncMap))
	assert.Equal(t, "four", HandleQuery("get 4", testSyncMap))
}

func TestHandleQuery__DelNonExisting__NotFound(t *testing.T){
	assert.Equal(t, "Value with this key not found", HandleQuery("del 5", testSyncMap))
}

func TestHandleQuery__DelExisting__Deleted(t *testing.T){
	assert.Equal(t, "three", HandleQuery("del 3", testSyncMap))
	assert.Equal(t, "Value with this key not found", HandleQuery("get 3", testSyncMap))
}

func TestHandleQuery__KeysExistingCorrectWildcard__Found(t *testing.T){
	querySet := HandleQuery("keys \"/*\"", testSyncMap)
	assert.True(t, "1,2,4" == querySet || "1,4,2" == querySet || "2,1,4" == querySet || "2,4,1" == querySet || "4,1,2" == querySet || "4,2,1" == querySet)
}

func TestHandleQuery__KeysNotExisting__NotFound(t *testing.T){
	assert.Equal(t, "Keys with this pattern not found", HandleQuery("keys \"123\"", testSyncMap))
}

func TestHandleQuery__KeysIncorrectWildcard__IncorrectPattern(t *testing.T){
	assert.Equal(t, "Pattern is incorrect", HandleQuery("keys \"*\"", testSyncMap))
}

func TestHandleQuery__123__UndefinedQuery(t *testing.T){
	assert.Equal(t, "Undefined query", HandleQuery("123", testSyncMap))
}
