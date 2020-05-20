package dhcp

import (
	"bytes"
	"dhcpickle/config"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func setCommonOptions(reply *dhcpv4.DHCPv4) {
	reply.ServerIPAddr = config.Config.ServerIdentifier
	reply.UpdateOption(dhcpv4.OptSubnetMask(config.Config.SubnetMask))                 // 1
	reply.UpdateOption(dhcpv4.OptRouter(config.Config.Router))                         // 3
	reply.UpdateOption(dhcpv4.OptDNS(config.Config.DNS...))                            // 6
	reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(config.Config.IPAddressLeaseTime)) // 51
	reply.UpdateOption(dhcpv4.OptServerIdentifier(config.Config.ServerIdentifier))     // 54
}

func getIP(hwAddr []byte) []byte {
	req, err := http.NewRequest("POST", config.Config.Endpoint, bytes.NewBuffer(hwAddr))
	if err != nil {
		log.Printf("error during http request: %v", err)
		return nil
	}
	authHeader := config.Config.AuthHeader
	authToken := config.Config.AuthToken
	if authHeader != "" && authToken != "" {
		req.Header.Add(authHeader, authToken)
	}
	resp, err := config.Client.Do(req)
	if err != nil {
		log.Printf("error during http response: %v", err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading body to string: %v", err)
		return nil
	}

	err = resp.Body.Close()
	if err != nil {
		log.Printf("error while closing body after read: %v", err)
	}

	return body
}

func DORAHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	log.Println(m.Summary())

	if m == nil {
		log.Println("Packet is nil!")
		return
	}
	if m.OpCode != dhcpv4.OpcodeBootRequest {
		log.Println("Not a BootRequest!")
		return
	}
	reply, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		log.Printf("NewReplyFromRequest failed: %v", err)
		return
	}

	// set all options on both offer and ack since clients may take either as the source of truth
	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))

		setCommonOptions(reply)
		ip := getIP(m.ClientHWAddr)
		if ip == nil {
			return
		}
		reply.YourIPAddr = ip
	case dhcpv4.MessageTypeRequest:
		ip := dhcpv4.GetIP(dhcpv4.OptionServerIdentifier, m.Options)
		if ip.Equal(config.Config.ServerIdentifier) {
			reply.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))

			setCommonOptions(reply)
			reply.YourIPAddr = dhcpv4.GetIP(dhcpv4.OptionRequestedIPAddress, m.Options) // 50
		} else {
			log.Printf("request was for: %v", ip)
			return
		}
	default:
		log.Printf("Unhandled message type: %v", mt)
		return
	}

	if _, err := conn.WriteTo(reply.ToBytes(), peer); err != nil {
		log.Printf("Cannot reply to client: %v", err)
	}
}
