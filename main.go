package main

import (
	"dhcpickle/config"
	"dhcpickle/dhcp"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"log"
	"net"
)

func main() {
	config.InitConfig()

	addr := &net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: dhcpv4.ServerPort,
	}
	server, err := server4.NewServer("", addr, dhcp.DORAHandler)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
