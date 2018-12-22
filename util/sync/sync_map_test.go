package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMap__SameMap__Success(t *testing.T) {
	testMap := map[string]string{"1": "one"}
	syncMap := Map{}
	syncMap.storage = &testMap
	assert.Equal(t, &testMap, syncMap.GetMap())
}

func TestGetMap__DiffMap__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one"}
	syncMap := Map{}
	syncMap.storage = &testMap
	diffMap := map[string]string{"2": "two"}
	assert.NotEqual(t, &diffMap, syncMap.GetMap())
}

func TestSetMap__SameMap__Success(t *testing.T) {
	testMap := map[string]string{"1": "one"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	assert.Equal(t, &testMap, syncMap.GetMap())
}

func TestSetMap__DiffMap__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	diffMap := map[string]string{"2": "two"}
	assert.NotEqual(t, &diffMap, syncMap.GetMap())
}

func TestGet__Existing__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	assert.Equal(t, "one", syncMap.Get("1"))
}

func TestGet__NotExisting__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	assert.Equal(t, "Value with this key not found", syncMap.Get("4"))
}

func TestSet__Changed__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	syncMap.Set("1", "123")
	assert.Equal(t, "123", syncMap.Get("1"))
}

func TestDel__DeletedAndReturned__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	assert.Equal(t, "one", syncMap.Del("1"))
}

func TestDel__DeletedAndNotFound__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	syncMap.Del("1")
	assert.Equal(t, "Value with this key not found", syncMap.Del("1"))
}

func TestKeys__FoundWildcard__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	querySet := syncMap.Keys("/*")
	assert.True(t, "1,2,3" == querySet || "1,3,2" == querySet || "2,1,3" == querySet || "2,3,1" == querySet || "3,2,1" == querySet || "3,1,2" == querySet)

}

func TestKeys__IncorrectWildcard__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	assert.Equal(t, "Pattern is incorrect", syncMap.Keys("*"))
}

func TestKeys__NotFound__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := Map{}
	syncMap.SetMap(&testMap)
	assert.Equal(t, "Keys with this pattern not found", syncMap.Keys("123"))
}
