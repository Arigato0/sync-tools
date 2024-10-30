package server

import "net"

func GetHostNames() ([]string, error) {
	return net.LookupAddr("example.com")
}
