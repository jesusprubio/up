package internal

import (
	"net"
	"testing"
)

// Tests the [AvailableInterfaces] function
func TestAvailableInterfaces(t *testing.T) {
	testCases := []struct {
		name               string
		interfaces         []net.Interface
		expectedInterfaces int
	}{
		{
			name: "Mixed Interfaces",
			interfaces: []net.Interface{
				{
					Name:  "lo",
					Flags: net.FlagLoopback | net.FlagUp,
				},
				{
					Name:  "eth0",
					Flags: 0,
				},
				{
					Name:  "docker0",
					Flags: net.FlagUp,
				},
				{
					Name:  "ens192",
					Flags: net.FlagUp,
				},
				{
					Name:  "ens224",
					Flags: net.FlagUp,
				},
			},
			expectedInterfaces: 2,
		},
		{
			name: "Multiple Active Interfaces",
			interfaces: []net.Interface{
				{
					Name:  "ens192",
					Flags: net.FlagUp,
				},
				{
					Name:  "eth0",
					Flags: net.FlagUp,
				},
			},
			expectedInterfaces: 2,
		},
		{
			name:               "No Interfaces",
			interfaces:         []net.Interface{},
			expectedInterfaces: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interfaces, err := AvailableInterfaces(tc.interfaces)
			if err != nil {
				if tc.expectedInterfaces == 0 && len(tc.interfaces) == 0 {
					return
				}
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(interfaces) != tc.expectedInterfaces {
				t.Errorf(
					"Expected %d interfaces, got %d",
					tc.expectedInterfaces,
					len(interfaces),
				)
			}

		})
	}
}
