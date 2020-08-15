package connman

import (
	"errors"
	"fmt"
	"time"
)

var (
	errNoNetworkAvailable error         = errors.New("no available network")
	errNoWifiInterface    error         = errors.New("no wifi interface")
	errNotConnected       error         = errors.New("not connected")
	configPath            string        = "/home/eco/wifi.conf"
	attempDelay           time.Duration = 2 * time.Second
)

//ConnectToNetwork main connect function
func ConnectToNetwork(name, pass string) error {
	var done bool = false
	var cmd string
	var attempt int = 0
	var attemptLimit int = 5
	err := interfaceUp()
	if err != nil {
		return err
	}
	for attempt < attemptLimit {
		ssid := ConnectedTo()
		if ssid == name {
			done = true
			break
		}
		if pass == "" {
			cmd = fmt.Sprintf(`sudo nmcli device wifi connect "%v"`, name)
		} else {
			cmd = fmt.Sprintf(`sudo nmcli device wifi connect "%v" password "%v"`, name, pass)
		}
		exe(cmd, "connect network")
		time.Sleep(attempDelay)
		attempt++
	}
	if done {
		return nil
	}
	return errNotConnected
}

//DisconnectFromNetwork main disconnect function
func DisconnectFromNetwork(name string) error {
	ssid := ConnectedTo()
	if ssid != name {
		return nil
	}
	cmd := fmt.Sprintf(`sudo nmcli con down id "%v"`, name)
	out, err := exe(cmd, "connect network")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	return nil
}
