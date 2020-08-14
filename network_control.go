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
	configPath            string = "/home/eco/wifi.conf"
	errNoNetworkAvailable error  = errors.New("no available network")
)

//GetNetworks get available wifi networks
func GetNetworks() ([]string, error) {
	cmd := "nmcli --fields SSID device wifi"
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

//HotspotOnWifiOff hotspot on wifi network off
func HotspotOnWifiOff(spot *Spot, wifi *WifiNetwork) error {
	var err error
	err = wifi.Disconnect()
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	err = spot.Up()
	if err != nil {
		return err
	}
	return nil
}

//HotspotOffWifiOn hotspot off wifi network on
func HotspotOffWifiOn(spot *Spot, wifi *WifiNetwork) error {
	var err error
	err = spot.Down()
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	err = wifi.Connect()
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
			configs = jin.MakeEmptyJson()
		} else {
			configs = penman.Read(configPath)
		}
	} else {
		configs = jin.MakeEmptyJson()
	}
	fmt.Println(configPath)
	fmt.Println(string(configs))
	return configs
}

//SaveNetwork save network to config path
func SaveNetwork(networkName, passphrase string) error {
	if result, _ := isThisNetworkSaved(networkName); result {
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

func isThisNetworkSaved(network string) (bool, error) {
	nets := ReadNetworks()
	_, err := jin.Get(nets, network)
	if err != nil {
		return false, err
	}
	return true, nil
}

//AvailableNetwork find saved available network
func AvailableNetwork() (string, string, error) {
	available, err := GetNetworks()
	if err != nil {
		return "", "", err
	}
	networks := ReadNetworks()
	list, err := jin.GetKeys(networks)
	if err != nil {
		return "", "", err
	}
	for _, a := range available {
		for _, n := range list {
			if a == n {
				pass, err := jin.GetString(networks, n)
				if err != nil {
					return "", "", err
				}
				return n, pass, nil
			}
		}
	}
	return "", "", errNoNetworkAvailable
}
