package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAndReadString__SameString__Success(t *testing.T) {
	file, _ := os.Create("test.txt")

	file.WriteString("123")

	file.Close()

	assert.Equal(t, "123", OpenAndReadString("test.txt"))

	os.Remove("test.txt")
}

func TestOpenAndReadString__OtherString__Fail(t *testing.T) {
	file, _ := os.Create("test.txt")

	file.WriteString("123")

	file.Close()

	assert.NotEqual(t, "124", OpenAndReadString("test.txt"))

	os.Remove("test.txt")
}

func TestOpenAndReadString__FileNotExist__Fail(t *testing.T) {
	assert.Panics(t, func() {
		OpenAndReadString("test.txt")
	})
}

func TestOpenAndRead__SameByteArray__Success(t *testing.T) {
	file, _ := os.Create("test.bin")

	file.Write([]byte("123"))

	file.Close()

	assert.Equal(t, []byte("123"), OpenAndRead("test.bin"))

	os.Remove("test.bin")
}

func TestOpenAndRead__DifferentByteArray__Success(t *testing.T) {
	file, _ := os.Create("test.bin")

	file.Write([]byte("123"))

	file.Close()

	assert.NotEqual(t, []byte("124"), OpenAndRead("test.bin"))

	os.Remove("test.bin")
}

func TestOpenAndRead__FileNotExist__Fail(t *testing.T) {
	assert.Panics(t, func() {
		OpenAndRead("test.bin")
	})
}
