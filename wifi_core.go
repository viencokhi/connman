package connman

import (
	"fmt"
	"net"
	"strings"
)

func isInterfaceUp() (bool, error) {
	cmd := fmt.Sprintf(`ip link show "%v"`, wifiInterface)
	out, _ := exe(cmd, "is interface up")
	if out == "" {
		return false, nil
	}
	start := strings.Index(out, "<")
	end := strings.Index(out, ">")
	out = out[start : end+1]
	up := strings.Contains(out, "UP")
	return up, nil
}

func interfaceUp() error {
	isUp, err := isInterfaceUp()
	if err != nil {
		return err
	}
	if isUp {
		return nil
	}
	cmd := fmt.Sprintf(`sudo ip link set "%v" up`, wifiInterface)
	out, err := exe(cmd, "interface up")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	return nil
}

func interfaceDown() error {
	isUp, err := isInterfaceUp()
	if err != nil {
		return err
	}
	if !isUp {
		return nil
	}
	cmd := fmt.Sprintf(`sudo ip link set "%v" down`, wifiInterface)
	out, err := exe(cmd, "interface down")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	return nil
}

func getWifiInterface() (string, error) {
	ints, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, nt := range ints {
		if strings.HasPrefix(nt.Name, "w") {
			return nt.Name, nil
		}
	}
	return "", errNoWifiInterface
}

//GetNetworks get available wifi networks
func GetNetworks() ([]string, error) {
	err := interfaceUp()
	if err != nil {
		return nil, err
	}
	cmd := "sudo nmcli --fields SSID device wifi"
	out, err := exe(cmd, "conn up")
	if err != nil {
		return nil, fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	networks := make([]string, 0, 8)
	tokens := strings.Split(string(out), "\n")
	tokens = tokens[1:]
	for _, t := range tokens {
		if t == "" {
			continue
		}
		networks = append(networks, strings.TrimSpace(t))
	}
	return networks, nil
}
