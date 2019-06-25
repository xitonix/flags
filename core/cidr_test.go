package core_test

import (
	"net"
	"testing"

	"go.xitonix.io/flags/core"
)

func TestCIDR_Equals(t *testing.T) {
	testCases := []struct {
		title    string
		c1       core.CIDR
		c2       core.CIDR
		expected bool
	}{
		{
			title:    "empty cidr inputs",
			expected: true,
		},
		{
			title:    "nil networks with unequal IP Addresses",
			expected: false,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.2"),
			},
		},
		{
			title:    "nil networks with equal IP Addresses",
			expected: true,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
			},
		},
		{
			title:    "first nil network with non nil empty second network",
			expected: true,
			c1: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: nil,
			},
			c2: core.CIDR{
				Network: &net.IPNet{},
				IP:      net.ParseIP("192.168.1.1"),
			},
		},
		{
			title:    "first nil network with non empty second network",
			expected: false,
			c1: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: nil,
			},
			c2: core.CIDR{
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask([]byte{0x1}),
				},
				IP: net.ParseIP("192.168.1.1"),
			},
		},
		{
			title:    "second nil network with non nil empty first network",
			expected: true,
			c1: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{},
			},
			c2: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: nil,
			},
		},
		{
			title:    "second nil network with non empty first network",
			expected: false,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask([]byte{0x1}),
				},
			},
			c2: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: nil,
			},
		},
		{
			title:    "empty networks with empty IP Addresses",
			expected: true,
			c1: core.CIDR{
				IP:      net.IP{},
				Network: &net.IPNet{},
			},
			c2: core.CIDR{
				IP:      net.IP{},
				Network: &net.IPNet{},
			},
		},
		{
			title:    "empty networks with equal IP Addresses",
			expected: true,
			c1: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{},
			},
			c2: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{},
			},
		},
		{
			title:    "non empty networks with unequal network IPs and nil masks",
			expected: false,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: nil,
				},
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.5"),
					Mask: nil,
				},
			},
		},
		{
			title:    "non empty networks with equal network IPs and nil masks",
			expected: true,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: nil,
				},
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: nil,
				},
			},
		},
		{
			title:    "non empty networks with equal network IPs and empty masks",
			expected: true,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask{},
				},
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask{},
				},
			},
		},
		{
			title:    "non empty networks with equal network IPs and unequal masks",
			expected: false,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask([]byte{0x1}),
				},
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask([]byte{0x2}),
				},
			},
		},
		{
			title:    "non empty networks with equal network IPs and equal masks",
			expected: true,
			c1: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask([]byte{0x1}),
				},
			},
			c2: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.4"),
					Mask: net.IPMask([]byte{0x1}),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := tc.c1.Equals(tc.c2)
			if actual != tc.expected {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}

func TestCIDR_String(t *testing.T) {
	ip, n, _ := net.ParseCIDR("192.168.1.1/24")
	invalidNetwork := &net.IPNet{
		IP:   net.ParseIP("192.168.1.2"),
		Mask: net.IPMask([]byte{1}),
	}
	testCases := []struct {
		title    string
		cidr     core.CIDR
		expected string
	}{
		{
			title:    "empty CIDR",
			expected: "",
		},
		{
			title: "non empty IP address and nil network",
			cidr: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: nil,
			},
			expected: "192.168.1.1",
		},
		{
			title: "non empty IP address and empty network",
			cidr: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{},
			},
			expected: "192.168.1.1",
		},
		{
			title: "non empty IP address and a network with nil mask",
			cidr: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.2"),
					Mask: nil,
				},
			},
			expected: "192.168.1.1",
		},
		{
			title: "non empty IP address and a network with empty mask",
			cidr: core.CIDR{
				IP: net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{
					IP:   net.ParseIP("192.168.1.2"),
					Mask: net.IPMask{},
				},
			},
			expected: "192.168.1.1",
		},
		{
			title: "non empty IP address and a none empty network",
			cidr: core.CIDR{
				IP: ip,
				Network: &net.IPNet{
					IP:   n.IP,
					Mask: n.Mask,
				},
			},
			expected: "192.168.1.1-192.168.1.0/24",
		},
		{
			title: "non empty IP address and invalid network",
			cidr: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: invalidNetwork,
			},
			expected: "192.168.1.1-" + invalidNetwork.String(),
		},
		{
			title: "empty IP address and invalid network",
			cidr: core.CIDR{
				Network: invalidNetwork,
			},
			expected: invalidNetwork.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := tc.cidr.String()
			if actual != tc.expected {
				t.Errorf("Expected: '%v', Actual: '%v'", tc.expected, actual)
			}
		})
	}
}

func TestCIDR_IsValid(t *testing.T) {
	ip, n, _ := net.ParseCIDR("192.168.1.1/24")
	testCases := []struct {
		title                 string
		cidr                  core.CIDR
		expected              bool
		expectedOriginalValue string
	}{
		{
			title:    "empty cidr",
			expected: false,
		},
		{
			title:    "cidr built without using the constructor",
			expected: false,
			cidr: core.CIDR{
				IP: ip,
				Network: &net.IPNet{
					IP:   n.IP,
					Mask: n.Mask,
				},
			},
		},
		{
			title:    "invalid cidr built using the constructor",
			expected: false,
			cidr:     core.NewCIDR("invalid"),
		},
		{
			title:    "invalid cidr built using an empty string",
			expected: false,
			cidr:     core.NewCIDR("   "),
		},
		{
			title:                 "valid cidr built using the constructor",
			expected:              true,
			cidr:                  core.NewCIDR("192.168.1.1/24"),
			expectedOriginalValue: "192.168.1.1/24",
		},

		{
			title:                 "valid cidr built using the constructor with white spaced input",
			expected:              true,
			cidr:                  core.NewCIDR("   192.168.1.1/24   "),
			expectedOriginalValue: "192.168.1.1/24",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := tc.cidr.IsValid()
			if actual != tc.expected {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}
			ao := tc.cidr.OriginalValue()
			if ao != tc.expectedOriginalValue {
				t.Errorf("Expected Original Value: '%v', Actual: '%v'", tc.expectedOriginalValue, ao)
			}
		})
	}
}
