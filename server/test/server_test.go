package test

import (
	"GoHomework/server/server"
	"strings"
	"testing"
)

const testDataFile = "data/test_data.json"

func TestServer(t *testing.T) {
	err := server.LoadData(testDataFile)
	if err != nil {
		t.Error("[ERR] LoadData failed. Cannot load from " + testDataFile + ".")
	}

	value := server.GetEntry("Key")
	if strings.Compare(value, "Value") != 0 {
		t.Error("[ERR] GetValue failed. Input: key=`Key`.")
	}

	value = server.GetEntry("ItDoesNotExist")
	if strings.Compare(value, "(nil)") != 0 {
		t.Error("[ERR] GetValue failed. Input: key=`ItDoesNotExist`.")
	}

	value = server.GetEntry("")
	if strings.Compare(value, "(nil)") != 0 {
		t.Error("[ERR] GetValue failed. Input: key=``.")
	}

	server.AddEntry("NewKey", "NewValue")
	removeResult := server.RemoveEntries([]string{"NewKey"})
	if strings.Compare(removeResult, "(integer) 1") != 0 {
		t.Error("[ERR] RemoveEntries failed. Input: keys={`NewKey1`}.")
	}

	server.AddEntry("NewKey1", "NewValue1")
	server.AddEntry("NewKey2", "NewValue2")
	removeResult = server.RemoveEntries([]string{"NewKey1", "NewKey2"})
	if strings.Compare(removeResult, "(integer) 2") != 0 {
		t.Error("[ERR] RemoveEntries failed. Input: keys={`NewKey1`, `NewKey2`}.")
	}

	keys := server.ShowAllKeys()
	if strings.Compare(keys, "\"Key\" ") != 0 {
		t.Error("[ERR] ShowAllKeys failed.")
	}

	server.AddEntry("NewKey1", "NewValue1")
	server.AddEntry("NewKey2", "NewValue2")
	server.AddEntry("NewKey3", "NewValue3")
	keys = server.FindKeys("NewKe*")
	if !strings.Contains(keys, "NewKey1") &&
		!strings.Contains(keys, "NewKey2") &&
		!strings.Contains(keys, "NewKey3") {
		t.Error("[ERR] FindKeys failed. Input: pattern=`NewKey*`.")
	}
	server.RemoveEntries([]string{"NewKey1", "NewKey2", "NewKey3"})

	server.SaveData(testDataFile)
}
