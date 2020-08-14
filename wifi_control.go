package connman

import (
	"errors"
	"fmt"
)

//WifiNetwork wifi connection struct
type WifiNetwork struct {
	name string
	pass string
	up   bool
}

var wiconnPath string = "/home/eco/go/src/connman/scripts/./wiconn.sh"

var errArgumentCount = errors.New("argument count error: you can pass network name and password if needed")

//NewWifiNetwork WifiNetwork constructor
func NewWifiNetwork(name, pass string) *WifiNetwork {
	return &WifiNetwork{name: name, pass: pass, up: false}
}

//Connect main network connection function
func (n *WifiNetwork) Connect() error {
	if n.up {
		return nil
	}
	var cmd string
	if n.pass == "" {
		cmd = fmt.Sprintf(`sudo "%v" connect "%v"`, wiconnPath, n.name)
	} else {
		cmd = fmt.Sprintf(`sudo "%v" connect "%v" "%v"`, wiconnPath, n.name, n.pass)
	}
	out, err := exe(cmd, "connect network")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	n.up = true
	return nil
}

//Disconnect main network connection function
func (n *WifiNetwork) Disconnect() error {
	if !n.up {
		return nil
	}
	var cmd string
	if n.pass == "" {
		cmd = fmt.Sprintf(`sudo "%v" disconnect "%v"`, wiconnPath, n.name)
	} else {
		cmd = fmt.Sprintf(`sudo "%v" disconnect "%v" "%v"`, wiconnPath, n.name, n.pass)
	}
	out, err := exe(cmd, "disconnect network")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	n.up = false
	return nil
}