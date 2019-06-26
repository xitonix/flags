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
			title:    "nil input inputs",
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
			actual := tc.c1.Equal(tc.c2)
			if actual != tc.expected {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}

func TestCIDR_FullString(t *testing.T) {
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
			cidr:     core.CIDR{},
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
			actual := tc.cidr.FullString()
			if actual != tc.expected {
				t.Errorf("Expected: '%v', Actual: '%v'", tc.expected, actual)
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
			cidr:     core.CIDR{},
			expected: "",
		},
		{
			title: "non empty IP address and nil network",
			cidr: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: nil,
			},
			expected: "",
		},
		{
			title: "non empty IP address and empty network",
			cidr: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: &net.IPNet{},
			},
			expected: "",
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
			expected: "",
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
			expected: "",
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
			expected: "192.168.1.1/24",
		},
		{
			title: "non empty IP address and invalid network",
			cidr: core.CIDR{
				IP:      net.ParseIP("192.168.1.1"),
				Network: invalidNetwork,
			},
			expected: "",
		},
		{
			title: "empty IP address and invalid network",
			cidr: core.CIDR{
				Network: invalidNetwork,
			},
			expected: "",
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

func TestParseCIDR(t *testing.T) {
	ip, n, _ := net.ParseCIDR("192.168.1.1/24")
	testCases := []struct {
		title              string
		input              string
		expected           *core.CIDR
		expectedString     string
		expectedFullString string
		expectErr          bool
	}{
		{
			title:              "empty input",
			expected:           nil,
			expectedString:     "",
			expectedFullString: "",
			expectErr:          true,
		},
		{
			title:              "white space input",
			expected:           nil,
			input:              "     ",
			expectedString:     "",
			expectedFullString: "",
			expectErr:          true,
		},
		{
			title: "valid input input with no white space",
			expected: &core.CIDR{
				IP:      ip,
				Network: n,
			},
			input:              "192.168.1.1/24",
			expectedFullString: "192.168.1.1-192.168.1.0/24",
			expectedString:     "192.168.1.1/24",
			expectErr:          false,
		},
		{
			title: "valid input input with white space sting",
			expected: &core.CIDR{
				IP:      ip,
				Network: n,
			},
			input:              "    192.168.1.1/24    ",
			expectedFullString: "192.168.1.1-192.168.1.0/24",
			expectedString:     "192.168.1.1/24",
			expectErr:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual, err := core.ParseCIDR(tc.input)
			if tc.expectErr && err == nil {
				t.Error("Expected to receive an error, but received nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Did not expect to receive an error, but received %s", err)
			}

			if err != nil {
				return
			}
			if !actual.Equal(*tc.expected) {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}

			actualFullString := actual.FullString()
			if actualFullString != tc.expectedFullString {
				t.Errorf("Expected full string: %v, Actual: %v", tc.expectedFullString, actualFullString)
			}

			actualString := actual.String()
			if actualString != tc.expectedString {
				t.Errorf("Expected string: %v, Actual: %v", tc.expectedString, actualString)
			}
		})
	}
}
