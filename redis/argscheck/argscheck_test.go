package argscheck

import "testing"

func TestStart(t *testing.T) {
	tables := []struct {
		args    []string
		gPort   string
		gMemory string
		gIP     string
		Port    string
		Memory  string
		IP      string
	}{
		{[]string{"1", "--mode", "disk", "--port", "8080", "-h", "127.0.0.7"}, ":9090", " ", "127.0.0.1", ":8080", "disk", "127.0.0.7"},
		{[]string{"2", "-m", "disk"}, ":9090", " ", "127.0.0.1", ":9090", "disk", "127.0.0.1"},
		{[]string{"3", "--port", "7070", "-h", "127.0.0.9"}, ":9090", " ", "127.0.0.1", ":7070", " ", "127.0.0.9"},
		{[]string{"4", "--port", "6060", "-h", "127.0.0.8"}, ":9090", " ", "127.0.0.1", ":6060", " ", "127.0.0.8"},
		{[]string{"5", "-h", "127.0.0.5"}, ":9090", " ", "127.0.0.1", ":9090", " ", "127.0.0.5"},
		{[]string{"6", "--port", "4040"}, ":9090", " ", "127.0.0.1", ":4040", " ", "127.0.0.1"},
	}
	for _, table := range tables {
		total, total2, total3 := Start(table.args, table.gPort, table.gMemory, table.gIP)
		if total != table.Port {
			t.Errorf("total %s, args %s, gPort %s, gMemory %s, gIP %s, Port %s, Memory %s, IP %s", total, table.args, table.gPort, table.gMemory, table.gIP, table.Port, table.Memory, table.IP)
		}
		if total2 != table.IP {
			t.Errorf("total2 %s, args %s, gPort %s, gMemory %s, gIP %s, Port %s, Memory %s, IP %s", total2, table.args, table.gPort, table.gMemory, table.gIP, table.Port, table.Memory, table.IP)
		}
		if total3 != table.Memory {
			t.Errorf("total3 %s, args %s, gPort %s, gMemory %s, gIP %s, Port %s, Memory %s, IP %s", total3, table.args, table.gPort, table.gMemory, table.gIP, table.Port, table.Memory, table.IP)
		}
	}
}
