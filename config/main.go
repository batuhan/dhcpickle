package config

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type config struct {
	SubnetMask         net.IPMask
	Router             net.IP
	DNS                []net.IP
	IPAddressLeaseTime time.Duration
	ServerIdentifier   net.IP
}

var Config = config{}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func copyAndAppend(slice []byte, elems ...byte) net.IP {
	newSlice := append(slice[:0:0], slice...)
	return append(newSlice, elems...)
}

func InitConfig() {
	outboundIP := getOutboundIP()
	envPrefix := "DHCP_"

	var prefix []byte
	prefixString := os.Getenv(envPrefix + "PREFIX")
	if prefixString == "" {
		prefix = outboundIP[:3]
	} else {
		prefix = []byte(prefixString)
	}

	subnetMask := os.Getenv(envPrefix + "SUBNET_MASK")
	if subnetMask == "" {
		Config.SubnetMask = net.IPMask{255, 255, 255, 0}
	} else {
		Config.SubnetMask = net.IPMask(net.ParseIP(subnetMask))
	}

	router := os.Getenv(envPrefix + "ROUTER")
	if router == "" {
		Config.Router = copyAndAppend(prefix, 1)
	} else {
		Config.Router = net.ParseIP(router)
	}

	dns := os.Getenv(envPrefix + "DNS")
	if dns == "" {
		Config.DNS = []net.IP{{8, 8, 8, 8}, {8, 8, 4, 4}}
	} else {
		var dnsResults []net.IP
		dnsStrings := strings.Split(dns, ",")
		for _, dnsString := range dnsStrings {
			dnsResults = append(dnsResults, net.ParseIP(dnsString))
		}
		Config.DNS = dnsResults
	}

	ipAddressLeaseTime := os.Getenv(envPrefix + "IP_ADDRESS_LEASE_TIME")
	if ipAddressLeaseTime == "" {
		Config.IPAddressLeaseTime = time.Hour * 24
	} else {
		var err error
		Config.IPAddressLeaseTime, err = time.ParseDuration(ipAddressLeaseTime)
		if err != nil {
			log.Fatalln(err)
		}
	}

	serverIdentifier := os.Getenv(envPrefix + "SERVER_IDENTIFIER")
	if serverIdentifier == "" {
		Config.ServerIdentifier = outboundIP
	} else {
		Config.ServerIdentifier = net.ParseIP(serverIdentifier)
	}

	log.Println("generated config", Config)
}
