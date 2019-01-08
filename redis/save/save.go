package save

import (
	"fmt"
	"os"
)

//SaveOnDisk - saves data into file on disk
func SaveOnDisk(info string) (int, error) {
	Path := "Info.txt"
	//Существующий файл с таким же именем будут перезаписан
	var fl, err = os.Create(Path)
	if err != nil {
		fmt.Printf("Error create: %v", err)
	}
	defer fl.Close()
	var ByteWrttn, errWrt = fl.WriteString(info)
	if errWrt != nil {
		fmt.Printf("Error write: %v", errWrt)
	}
	fmt.Printf("Info.txt written: %v\n", ByteWrttn)
	return ByteWrttn, nil
}
