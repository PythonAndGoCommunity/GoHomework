package collection

import (
	"net"
)

func Index(slice []interface{}, value interface{}) int{
	for index, slice_value := range slice {
		if slice_value == value {
			return index
		}
	}
	return -1
}

func ConnIndex(slice []net.Conn, value net.Conn) int{
	for index, slice_value := range slice {
		if slice_value == value {
			return index
		}
	}
	return -1
}