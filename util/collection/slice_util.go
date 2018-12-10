package collection

import (
	"net"
)

// ConnIndex returns index of neccessary element in net.Conn slice if found, otherwise -1.
func ConnIndex(slice []net.Conn, value net.Conn) int{
	for index, slice_value := range slice {
		if slice_value == value {
			return index
		}
	}
	return -1
}