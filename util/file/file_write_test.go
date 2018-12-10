package file

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestCreateAndWriteString__SameString__Success(t *testing.T){
	CreateAndWriteString("test.txt", "123")
	assert.Equal(t, "123", OpenAndReadString("test.txt"))
	os.Remove("test.txt")
}

func TestCreateAndWriteString__OtherString__Fail(t *testing.T){
	CreateAndWriteString("test.txt", "123")
	assert.NotEqual(t, "124", OpenAndReadString("test.txt"))
	os.Remove("test.txt")
}

func TestCreateAndWrite__SameByteArray__Success(t *testing.T){
	CreateAndWrite("test.bin", []byte("123"))
	assert.Equal(t, []byte("123"), OpenAndRead("test.bin"))
	os.Remove("test.bin")
}

func TestCreateAndWrite__OtherByteArray__Success(t *testing.T){
	CreateAndWrite("test.bin", []byte("123"))
	assert.NotEqual(t, []byte("124"), OpenAndRead("test.bin"))
	os.Remove("test.bin")
}


