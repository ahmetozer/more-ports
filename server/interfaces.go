package server

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// listUpInterfaces List all up interfaces
func listUpInterfaces() ([]net.Interface, error) {
	allInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var upInterfaces []net.Interface
	for i := range allInterfaces {
		//https://golang.org/src/net/interface.go#40
		var a = allInterfaces[i].Flags & (1 << uint(0))
		if a.String() == "up" {
			upInterfaces = append(upInterfaces, allInterfaces[i])
		}

	}
	return upInterfaces, nil
}

//DefaultRoute Get default IPv4 route with 192.88.99.1
func DefaultRoute() (string, error) {
	if defaultRoute, err := netlink.RouteGet(net.IP{192, 88, 99, 1}); err == nil {
		return fmt.Sprintf("%s", defaultRoute[0].Gw), nil
	} else {
		return "", err
	}

}
