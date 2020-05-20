package dhcp

import (
	"dhcpickle/config"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"log"
	"net"
)

func setCommonOptions(reply *dhcpv4.DHCPv4) {
	reply.ServerIPAddr = config.Config.ServerIdentifier
	reply.UpdateOption(dhcpv4.OptSubnetMask(config.Config.SubnetMask))                 // 1
	reply.UpdateOption(dhcpv4.OptRouter(config.Config.Router))                         // 3
	reply.UpdateOption(dhcpv4.OptDNS(config.Config.DNS...))                            // 6
	reply.UpdateOption(dhcpv4.OptIPAddressLeaseTime(config.Config.IPAddressLeaseTime)) // 51
	reply.UpdateOption(dhcpv4.OptServerIdentifier(config.Config.ServerIdentifier))     // 54
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
		reply.YourIPAddr = net.IP{185, 226, 95, 104}
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
