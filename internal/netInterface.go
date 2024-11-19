package internal

import (
	"fmt"
	"net"
	"strings"
)

type InterfaceInfo struct {
	Name string
	IP   net.IP
}

// AvailableInterfaces returns a list of active network interfaces
// with their IPv4 addresses that can connect to the outside world
func AvailableInterfaces(interfaces []net.Interface) ([]InterfaceInfo, error) {
	if len(interfaces) == 0 {
		return nil, fmt.Errorf("no interfaces provided")
	}
	var activeInterfaces []InterfaceInfo
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 ||
			iface.Flags&net.FlagUp == 0 ||
			isVirtualInterface(iface.Name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf(
				"Warning: Error fetching addresses for interface %s: %v\n",
				iface.Name,
				err,
			)
			continue
		}

		var ipV4 net.IP

		for _, addr := range addrs {
			parsedIP, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			if parsedIP.To4() != nil && !parsedIP.IsLoopback() {
				ipV4 = parsedIP
				break
			}
		}

		activeInterfaces = append(activeInterfaces, InterfaceInfo{
			Name: iface.Name,
			IP:   ipV4,
		})
	}

	return activeInterfaces, nil
}

func isVirtualInterface(name string) bool {
	virtualPrefixes := []string{
		"virbr", "vnet", "docker", "br-", "tun", "tap",
		"vmnet", "veth", "vbox", "wg",
		"kube", "cali", "flannel",
		"vmx", "vlan", "bond", "teredo",
	}
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
