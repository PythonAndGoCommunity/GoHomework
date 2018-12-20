package db

import "testing"

func newTestDb(m map[string]string) *DataBase {
	db := DataBase{mode: "memory", m: m}
	return &db
}

func TestDataBase(t *testing.T) {
	t.Log("\tTest DataBase functions")

	// m := map[string]string{"a": "1", "b": "2"}
	// db := DataBase{m: m}

	type testTuple struct {
		key   string
		value string
	}

	t.Run("Get", func(t *testing.T) {
		tests := []testTuple{
			{"a", "1"},
			{"b", "2"},
			{"c", ""},
		}
		db := newTestDb(map[string]string{"a": "1", "b": "2"})
		for _, tt := range tests {
			actual, _ := db.Get(tt.key)
			if actual != tt.value {
				t.Errorf("function Get(%s): expected %s, actual %s", tt.key, tt.value, actual)
			}
		}
		key := "c"
		expected := false
		_, actual := db.Get(key)
		if actual != expected {
			t.Errorf("function Get(%s): expected %v, actual %v", key, expected, actual)
		}
	})

	t.Run("Set", func(t *testing.T) {
		tests := []testTuple{
			{"a", "1"},
			{"b", "2"},
			{"c", "3"},
		}
		db := newTestDb(map[string]string{})
		for _, tt := range tests {
			db.Set(tt.key, tt.value)
			actual, _ := db.Get(tt.key)
			if actual != tt.value {
				t.Errorf("function Get(%v): expected %v, actual %v", tt.key, tt.value, actual)
			}
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tests := map[string]string{"a": "1", "b": "2"}
		db := newTestDb(tests)
		tests["c"] = "3" // for a key not presented in test db
		for k := range tests {
			db.Delete(k)
			_, has := db.Get(k)
			expected := false
			if has != expected {
				t.Errorf("function Delete(%s): expected %v, actual %v", k, expected, has)
			}
		}
	})
}
