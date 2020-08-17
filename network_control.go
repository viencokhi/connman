package connman

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	errNoNetworkAvailable error         = errors.New("no available network")
	errNoWifiInterface    error         = errors.New("no wifi interface")
	errNotConnected       error         = errors.New("not connected")
	attempDelay           time.Duration = 3 * time.Second
)

func getIP() string {
	ints, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, in := range ints {
		prefix := in.Name[0]
		if prefix == 'l' {
			continue
		}
		addr, err := in.Addrs()
		if err != nil {
			continue
		}
		if addr == nil {
			continue
		}
		ip := addr[0].String()
		tokens := strings.Split(ip, "/")
		return tokens[0]
	}
	return ""
}

func getConnections() map[string]string {
	connections := make(map[string]string)
	connections["wifi"] = "false"
	connections["ethernet"] = "false"
	ints, err := net.Interfaces()
	if err != nil {
		return connections
	}
	for _, in := range ints {
		addr, err := in.Addrs()
		if err != nil {
			continue
		}
		if addr == nil {
			continue
		}
		prefix := in.Name[0]
		switch prefix {
		case 'l':
			continue
		case 'w':
			connections["wifi"] = "true"
			break
		case 'e':
			connections["ethernet"] = "true"
			break
		}
	}
	return connections
}

//ConnectedToInternet internet connection check
func connectedToInternet() bool {
	resp, err := http.Get("https://www.google.com/")
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func disconnectFromAllConnection() error {
	cmd := `sudo nmcli -f NAME con sh`
	out, err := exe(cmd, "get connection")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	tokens := strings.Split(out, "\n")
	tokens = tokens[1:]
	if len(tokens) == 0 {
		return nil
	}
	for _, t := range tokens {
		con := strings.TrimSpace(t)
		cmd = fmt.Sprintf(`sudo nmcli con down "%v"`, con)
		out, err = exe(cmd, "down connection")
		if err != nil {
			return fmt.Errorf("error:%v, out:%v", err.Error(), out)
		}
	}
	return nil
}

//ConnectToNetwork main connect function
func connectToNetwork(name, pass string) error {
	var done bool = false
	var cmd string
	var attempt int = 0
	var attemptLimit int = 5
	err := interfaceUp()
	if err != nil {
		return err
	}
	for attempt < attemptLimit {
		ssid := connectedTo()
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
func disconnectFromNetwork(name string) error {
	ssid := connectedTo()
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
