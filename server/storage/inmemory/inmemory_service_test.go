package inmemory

import (
	"NonRelDB/util/file"
	"NonRelDB/util/json"
	"NonRelDB/util/sync"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStorage__Empty__Fail(t *testing.T) {
	var syncMap *sync.Map
	assert.Equal(t, syncMap, GetStorage())
}

func TestGetStorage__Existing__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := sync.Map{}
	syncMap.SetMap(&testMap)
	storage = &syncMap
	assert.Equal(t, &syncMap, GetStorage())
}

func TestSetStorage__NonEmpty__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := sync.Map{}
	syncMap.SetMap(&testMap)
	SetStorage(&syncMap)
	assert.Equal(t, &syncMap, GetStorage())
}

func TestSetStorage__Empty__Fail(t *testing.T) {
	emptyMap := sync.Map{}
	SetStorage(&emptyMap)
	assert.Equal(t, &emptyMap, GetStorage())
}

func TestInitDBInMemory__Success(t *testing.T) {
	assert.NotPanics(t, func() {
		InitDBInMemory()
	})
}

func TestInitDBFromStorage__Existing__Success(t *testing.T) {
	file.CreateAndWriteString("test.json", "{\"123\":\"123\"}")
	InitDBFromStorage("test.json")
	testMap := map[string]string{"123": "123"}
	assert.Equal(t, testMap, *(GetStorage().GetMap()))

	os.Remove("test.json")
}

func TestInitDBFromStorage__NotExisting__Success(t *testing.T) {
	InitDBFromStorage("test.json")
	testMap := map[string]string{}
	assert.Equal(t, testMap, *(GetStorage().GetMap()))

	os.Remove("test.json")
}

func TestRestoreDBFromDump__NonEmpty__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	RestoreDBFromDump(json.PackMapToJSON(testMap))
	assert.Equal(t, testMap, *(GetStorage().GetMap()))
}

func TestSaveDBToStorage__NonEmpty__Success(t *testing.T) {
	testMap := map[string]string{"1": "one", "2": "two", "3": "three"}
	syncMap := sync.Map{}
	syncMap.SetMap(&testMap)
	SetStorage(&syncMap)
	SaveDBToStorage("test.json")
	assert.Equal(t, "{\"1\":\"one\",\"2\":\"two\",\"3\":\"three\"}", file.OpenAndReadString("test.json"))

	os.Remove("test.json")
}
