package core

import (
	"net"
	"strings"
)

// CIDR represents a pair of an IP address and a Network to which the IP belong.
//
// For example, "192.0.2.1/24" has the IP address 192.0.2.1 and the network 192.0.2.0/24.
//
// You MUST only build a new CIDR object using NewCIDR function, otherwise it will be ignored by the library.
type CIDR struct {
	// IP address.
	IP net.IP
	// Network the network implied by the IP and prefix length.
	Network       *net.IPNet
	originalValue string
}

// NewCIDR creates a new CIDR object.
//
// You MUST only build a new CIDR object using this function, otherwise it will be ignored by the library.
// Also, if the input string does not follow RFC 4632 and RFC 4291 standards (i.e. "192.0.2.0/24" or "2001:db8::/32"),
// the result CIDR will be invalid and cannot be used by the library.
//
// This is an alias for `net.ParseCIDR(cidr)` to simplify the flags API.
func NewCIDR(cidr string) CIDR {
	cidr = strings.TrimSpace(cidr)
	i, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return CIDR{}
	}
	return CIDR{
		IP:            i,
		Network:       n,
		originalValue: cidr,
	}
}

// IsValid returns true if the CIDR has been successfully created using NewCIDR() function.
func (c CIDR) IsValid() bool {
	return len(c.originalValue) > 0
}

// OriginalValue returns the original string of the CIDR if it has been successfully created using NewCIDR() function.
func (c CIDR) OriginalValue() string {
	return c.originalValue
}

// Equals returns true if the CIDRs are equal.
func (c CIDR) Equals(to CIDR) bool {
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

// String returns the string representation of the CIDR in IP-Network/Length format.
func (c CIDR) String() string {
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

func isEmptyNet(n *net.IPNet) bool {
	if n == nil {
		return true
	}
	return len(n.IP) == 0 && len(n.Mask) == 0
}
