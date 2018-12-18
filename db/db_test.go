package db

import "testing"

func TestSet(t *testing.T) {
	t.Log("\tTest DataBase functions")

	m := map[string]string{"a": "1", "b": "2"}
	db := DataBase{m: m}

	type testTuple struct {
		key   string
		value string
	}

	t.Run("Delete", func(t *testing.T) {
		key := "b"
		expected := ""
		db.Delete(key)
		actual, _ := db.Get(key)
		if actual != expected {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	})

	t.Run("Get", func(t *testing.T) {
		tests := []testTuple{
			{"a", "1"},
			{"b", ""},
			{"c", ""},
		}
		for _, tt := range tests {
			actual, _ := db.Get(tt.key)
			if actual != tt.value {
				t.Errorf("function Get(%v): expected %v, actual %v", tt.key, tt.value, actual)

			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		tests := []testTuple{
			{"a", "1"},
			{"b", "2"},
		}
		for _, tt := range tests {
			db.Set(tt.key, tt.value)
			actual, _ := db.Get(tt.key)
			if actual != tt.value {
				t.Errorf("function Get(%v): expected %v, actual %v", tt.key, tt.value, actual)

			}
		}
	})

}
