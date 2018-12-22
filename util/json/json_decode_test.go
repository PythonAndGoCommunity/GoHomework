package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackMapFromJSON__SameMap__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	mapBytes, _ := json.Marshal(testMap)
	assert.Equal(t, testMap, *UnpackFromJSON(mapBytes))
}

func TestUnpackMapFromJSON__DiffMap__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	mapBytes, _ := json.Marshal(testMap)
	testMap["1"] = "1"
	assert.NotEqual(t, testMap, *UnpackFromJSON(mapBytes))
}

func TestUnpackMapFromJSON__SameMapWithIndent__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	mapBytes, _ := json.MarshalIndent(testMap, "", " ")
	assert.Equal(t, testMap, *UnpackFromJSON(mapBytes))
}

func TestUnpackMapFromJSON__DiffMapWithIndent__Fail(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	mapBytes, _ := json.MarshalIndent(testMap, "", " ")
	testMap["1"] = "1"
	assert.NotEqual(t, testMap, *UnpackFromJSON(mapBytes))
}
