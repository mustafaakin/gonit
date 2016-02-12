package main

import (
	"fmt"
	"net"
)

func main() {
	println("I will do networking and stuff later")
	fmt.Printf("WAT")

	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	} else {
		for _, iface := range ifaces {
			println(iface.Name, iface.HardwareAddr.String(), iface.MTU)
		}
	}
}
