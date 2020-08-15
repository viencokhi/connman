package connman

import (
	"fmt"
	"jin"
	"penman"
	"strings"
	"time"
)

var (
	configPath string = "/etc/wifi.conf"
)

//Status network and connection status
func Status() []byte {
	status := []byte(`{"internet":null, "connections":null, "connectedTo":null, "availableNetworks":null, "savedNetworks":null,"error":null}`)
	connections := jin.MakeJsonWithMap(getConnections())
	available, err := GetNetworks()
	if err != nil {
		status, _ = jin.AddKeyValueString(status, "error", err.Error())
		return status
	}
	saved := ReadNetworks()
	availableNetworks := jin.MakeArrayString(available)
	status, err = jin.SetBool(status, ConnectedToInternet(), "internet")
	if err != nil {
		status, _ = jin.SetString(status, err.Error(), "error")
		return status
	}
	status, err = jin.Set(status, connections, "connections")
	if err != nil {
		status, _ = jin.SetString(status, err.Error(), "error")
		return status
	}
	connected := ConnectedTo()
	if connected != "" {
		status, err = jin.SetString(status, connected, "connectedTo")
	}
	if err != nil {
		status, _ = jin.SetString(status, err.Error(), "error")
		fmt.Println("here", status)
		return status
	}
	status, err = jin.Set(status, availableNetworks, "availableNetworks")
	if err != nil {
		status, _ = jin.SetString(status, err.Error(), "error")
		return status
	}
	status, err = jin.Set(status, saved, "savedNetworks")
	if err != nil {
		status, _ = jin.SetString(status, err.Error(), "error")
		return status
	}
	status, err = jin.SetString(status, "null", "error")
	if err != nil {
		status, _ = jin.SetString(status, err.Error(), "error")
		return status
	}
	return status
}

//ConnectedTo returns SSID of wifi connection
func ConnectedTo() string {
	err := interfaceUp()
	if err != nil {
		return ""
	}
	cmd := fmt.Sprintf(`nmcli connection show | grep "%v"`, wifiInterface)
	out, _ := exe(cmd, "connected to")
	index := strings.Index(out, " ")
	if index == -1 {
		return ""
	}
	return out[:index]
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

//RemoveNetwork remove network from config path
func RemoveNetwork(network string) error {
	if result, _ := hasThisNetwork(network); !result {
		return nil
	}
	var err error
	configs := ReadNetworks()
	configs, err = jin.Delete(configs, network)
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
