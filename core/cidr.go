package core

import (
	"fmt"
	"net"
	"strings"
)

// CIDR represents a pair of an IP address and a Network to which the IP belong.
//
// For example, "192.0.2.1/24" has the IP address 192.0.2.1 and the network 192.0.2.0/24.
type CIDR struct {
	// IP address.
	IP net.IP
	// Network the network implied by the IP and prefix length.
	Network *net.IPNet
}

// ParseCIDR creates a new CIDR object.
//
// This is an alias for `net.ParseCIDR(input)`.
// The input string must be compliant with RFC 4632 and RFC 4291 standards (i.e. "192.0.2.0/24" or "2001:db8::/32"),
func ParseCIDR(cidr string) (*CIDR, error) {
	cidr = strings.TrimSpace(cidr)
	i, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &CIDR{
		IP:      i,
		Network: n,
	}, nil
}

// String returns the string representation of the original CIDR string in IP/Length format.
func (c CIDR) String() string {
	if len(c.IP) > 0 && c.Network != nil {
		ns := c.Network.String()
		idx := strings.Index(ns, "/")
		if idx > 0 {
			return fmt.Sprintf("%v%v", c.IP, ns[idx:])
		}
	}
	return ""
}

// FullString returns the string representation of the CIDR in IP-Network/Length format.
func (c CIDR) FullString() string {
	var buf string
	if len(c.IP) != 0 {
		buf += c.IP.String()
	}
	if c.Network != nil && len(c.Network.IP) > 0 && len(c.Network.Mask) > 0 {
		if len(buf) > 0 {
			buf += "-"
		}
		buf += c.Network.String()
	}
	return buf
}

// Equal returns true if the CIDRs are equal.
func (c CIDR) Equal(to CIDR) bool {
	if !c.IP.Equal(to.IP) {
		return false
	}

	if isEmptyNet(c.Network) {
		return isEmptyNet(to.Network)
	}

	if isEmptyNet(to.Network) {
		return isEmptyNet(c.Network)
	}

	if !c.Network.IP.Equal(to.Network.IP) || len(c.Network.Mask) != len(to.Network.Mask) {
		return false
	}

	for i, cb := range c.Network.Mask {
		if to.Network.Mask[i] != cb {
			return false
		}
	}

	return true
}

func isEmptyNet(n *net.IPNet) bool {
	if n == nil {
		return true
	}
	return len(n.IP) == 0 && len(n.Mask) == 0
}
