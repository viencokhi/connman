package connman

import (
	"errors"
	"fmt"
	"jin"
	"penman"
	"strings"
	"time"
)

var (
	configPath            string        = "/home/eco/wifi.conf"
	errNoNetworkAvailable error         = errors.New("no available network")
	attempDelay           time.Duration = 2 * time.Second
)

//GetNetworks get available wifi networks
func GetNetworks() ([]string, error) {
	cmd := "sudo nmcli --fields SSID device wifi"
	out, err := exe(cmd, "conn up")
	if err != nil {
		return nil, fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	networks := make([]string, 0, 8)
	tokens := strings.Split(string(out), "\n")
	tokens = tokens[1:]
	for _, t := range tokens {
		networks = append(networks, strings.TrimSpace(t))
	}
	return networks, nil
}

//HotspotOn hotspot on wifi network off
func HotspotOn(spot *Spot) error {
	err := spot.Up()
	if err != nil {
		return err
	}
	return nil
}

//ConnectAvailable scan, find and connect a saved available wifi network
func ConnectAvailable() error {
	wifi, err := AvailableNetwork()
	if err != nil {
		return err
	}
	if wifi == nil {
		return errNoNetworkAvailable
	}
	err = wifi.Connect()
	if err != nil {
		return err
	}
	return nil
}

//HotspotOff hotspot off wifi network on
func HotspotOff(spot *Spot) error {
	var err error
	err = spot.Down()
	if err != nil {
		return err
	}
	err = ConnectAvailable()
	if err != nil {
		return err
	}
	return nil
}

//ReadNetworks Read saved networks from config path
func ReadNetworks() []byte {
	var configs []byte
	if penman.IsFileExist(configPath) {
		if penman.IsFileEmpty(configPath) {
			return []byte(`{}`)
		}
		configs = penman.Read(configPath)
	} else {
		return []byte(`{}`)
	}
	return configs
}

//ReadNetwork Reads a saved network
func ReadNetwork(networkName string) (*WifiNetwork, error) {
	networks := ReadNetworks()
	pass, err := jin.GetString(networks, networkName)
	if err != nil {
		return nil, err
	}
	return &WifiNetwork{name: networkName, pass: pass, up: false}, nil
}

//SaveNetwork save network to config path
func SaveNetwork(networkName, passphrase string) error {
	if result, _ := hasThisNetwork(networkName); result {
		return nil
	}
	var err error
	var configs []byte
	if penman.IsFileExist(configPath) {
		if penman.IsFileEmpty(configPath) {
			configs = jin.MakeEmptyJson()
		} else {
			configs = penman.Read(configPath)
		}
	} else {
		configs = jin.MakeEmptyJson()
	}
	configs, err = jin.AddKeyValueString(configs, networkName, passphrase)
	if err != nil {
		return err
	}
	penman.OWrite(configPath, configs)
	return nil
}

func hasThisNetwork(network string) (bool, error) {
	nets := ReadNetworks()
	_, err := jin.Get(nets, network)
	if err != nil {
		return false, err
	}
	return true, nil
}

//AvailableNetwork find saved available network
func AvailableNetwork() (*WifiNetwork, error) {
	var done bool = false
	var attempt int = 0
	var attemptLimit int = 5
	for attempt < attemptLimit {
		if done {
			break
		}
		available, err := GetNetworks()
		if err != nil {
			return nil, err
		}
		networks := ReadNetworks()
		list, err := jin.GetKeys(networks)
		if err != nil {
			return nil, err
		}
		for _, a := range available {
			for _, n := range list {
				if a == n {
					pass, err := jin.GetString(networks, n)
					if err != nil {
						return nil, err
					}
					done = true
					return &WifiNetwork{name: n, pass: pass, up: false}, nil
				}
			}
		}
		time.Sleep(attempDelay)
		attempt++
	}
	return nil, errNoNetworkAvailable
}
