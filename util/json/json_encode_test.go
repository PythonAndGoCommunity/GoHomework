package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
)

func TestPackToJSON__SameMap__Success(t *testing.T){
	testMap := map[string]string{"1" : "one"}
	expBytes, _ := json.Marshal(testMap)
	assert.Equal(t, expBytes, PackToJSON("1", testMap["1"]))
}

func TestPackToJSON__OtherMap__Fail(t *testing.T){
	testMap := map[string]string{"1" : "one"}
	expBytes, _ := json.Marshal(testMap)
	assert.NotEqual(t, expBytes, PackToJSON("1", "two"))
}

func TestPackToJSONIndent__SameMap__Success(t *testing.T){
	testMap := map[string]string{"1" : "one"}
	expBytes, _ := json.MarshalIndent(testMap, "", " ")
	assert.Equal(t, expBytes, PackToJSONIndent("1", testMap["1"]))
}

func TestPackToJSONIndent__OtherMap__Fail(t *testing.T){
	testMap := map[string]string{"1" : "one"}
	expBytes, _ := json.MarshalIndent(testMap,"", " ")
	assert.NotEqual(t, expBytes, PackToJSONIndent("1", "two"))
}


func TestPackMapToJSON__SameMap__Success(t *testing.T){
	testMap := map[string]string{"1" : "one", "2" : "two", "3" : "three"}
	expBytes, _ := json.Marshal(testMap)
	assert.Equal(t, expBytes, PackMapToJSON(testMap))
}

func TestPackMapToJSON__DiffMap__Fail(t *testing.T){
	testMap := map[string]string{"1" : "one", "2" : "two", "3" : "three"}
	expBytes, _ := json.Marshal(testMap)
	testMap["1"] = "1"
	assert.NotEqual(t, expBytes, PackMapToJSON(testMap))
}

func TestPackMapToJSONIndent__SameMap__Success(t *testing.T){
	testMap := map[string]string{"1" : "one", "2" : "two", "3" : "three"}
	expBytes, _ := json.MarshalIndent(testMap, "", " ")
	assert.Equal(t, expBytes, PackMapToJSONIndent(testMap))
}

func TestPackMapToJSONIndent__DiffMap__Fail(t *testing.T){
	testMap := map[string]string{"1" : "one", "2" : "two", "3" : "three"}
	expBytes, _ := json.MarshalIndent(testMap, "", " ")
	testMap["1"] = "1"
	assert.NotEqual(t, expBytes, PackMapToJSONIndent(testMap))
}