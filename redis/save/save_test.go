package save

import (
	"testing"
)

func TestSaveOnDisk(t *testing.T) {
	tables := []struct {
		info    string
		byte    int
		testerr error
	}{
		{"hello world i am here stas here everyone is here", 48, nil},
	}
	for _, table := range tables {
		total, totalerr := SaveOnDisk(table.info)
		if total != table.byte {
			t.Errorf("total %v, byte %v, info %s", total, table.byte, table.info)
		}
		if totalerr != table.testerr {
			t.Errorf("totalerr %v, testerr %v, info %s", totalerr, table.testerr, table.info)
		}
	}
}
