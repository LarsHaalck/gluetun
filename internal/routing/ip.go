package routing

import (
	"errors"
	"fmt"
	"net"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

func IPIsPrivate(ip netip.Addr) bool {
	return ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

var (
	errInterfaceIPNotFound = errors.New("IP address not found for interface")
)

func ipMatchesFamily(ip netip.Addr, family int) bool {
	return (family == netlink.FAMILY_V6 && ip.Is6()) ||
		(family == netlink.FAMILY_V4 && (ip.Is4() || ip.Is4In6()))
}

func (r *Routing) assignedIP(interfaceName string, family int) (ip netip.Addr, err error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return ip, fmt.Errorf("network interface %s not found: %w", interfaceName, err)
	}
	addresses, err := iface.Addrs()
	if err != nil {
		return ip, fmt.Errorf("listing interface %s addresses: %w", interfaceName, err)
	}
	for _, address := range addresses {
		switch value := address.(type) {
		case *net.IPAddr:
			ip = netIPToNetipAddress(value.IP)
		case *net.IPNet:
			ip = netIPToNetipAddress(value.IP)
		default:
			continue
		}

		if !ipMatchesFamily(ip, family) {
			continue
		}

		// Ensure we don't return an IPv6-wrapped IPv4 address
		// since netip.Address String method works differently than
		// net.IP String method for this kind of addresses.
		return ip.Unmap(), nil
	}
	return ip, fmt.Errorf("%w: interface %s in %d addresses",
		errInterfaceIPNotFound, interfaceName, len(addresses))
}
