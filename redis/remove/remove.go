package remove

import "fmt"

//Remove - remove element from slice
func Remove(list []string, num int) []string {
	fmt.Printf("'remove' done\n")
	return append(list[:num], list[num+1:]...)
}
