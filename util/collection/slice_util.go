package collection

import (
	"net"
)

// Index returns index of neccessary element in interface{} slice if found, otherwise -1.
func Index(slice []interface{}, value interface{}) int{
	for index, slice_value := range slice {
		if slice_value == value {
			return index
		}
	}
	return -1
}

// ConnIndex returns index of neccessary element in net.Conn slice if found, otherwise -1.
func ConnIndex(slice []net.Conn, value net.Conn) int{
	for index, slice_value := range slice {
		if slice_value == value {
			return index
		}
	}
	return -1
}