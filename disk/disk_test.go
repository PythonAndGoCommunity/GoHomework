package disk

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	redislight "github.com/ITandElectronics/GoHomework"
)

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "redis-light-tests")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestNew(t *testing.T) {
	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	_, err := New(fmt.Sprintf("%s/test_db.json", dir))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDiskGet(t *testing.T) {
	disk, err := New("./fixtures/db.json")
	if err != nil {
		t.Fatal(err)
	}
	_, err = disk.Get("not_existing_key")
	if err == nil {
		t.Fatal("should return error")
	}
	if err != redislight.ErrKeyIsNotExists {
		t.Fatalf("%v - is not expected", err)
	}
	val, err := disk.Get("exists")
	if err != nil {
		t.Fatal(err)
	}
	if val == "" {
		t.Fatal("should return some value")
	}

}

func TestDiskSet(t *testing.T) {
	disk, err := New("./fixtures/db.json")
	if err != nil {
		t.Fatal(err)
	}
	err = disk.Set("new", "world")
	if err != nil {
		t.Fatal(err)
	}
	val, err := disk.Get("new")
	if err != nil {
		t.Fatal(err)
	}
	if val != "world" {
		t.Fatal("Set is not working")
	}

}

func TestDiskDel(t *testing.T) {

}

func TestDiskSync(t *testing.T) {

}
