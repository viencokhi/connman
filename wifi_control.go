package connman

//WifiNetwork wifi connection struct
type WifiNetwork struct {
	name string
	pass string
	up   bool
}

//NewWifiNetwork WifiNetwork constructor
func NewWifiNetwork(name, pass string) *WifiNetwork {
	return &WifiNetwork{name: name, pass: pass, up: false}
}

//Connect main network connection function
func (n *WifiNetwork) Connect() error {
	if n.up {
		return nil
	}
	err := connectToNetwork(n.name, n.pass)
	if err != nil {
		return err
	}
	n.up = true
	return nil
}

//Disconnect main network connection function
func (n *WifiNetwork) Disconnect() error {
	if !n.up {
		return nil
	}
	err := disconnectFromNetwork(n.name)
	if err != nil {
		return err
	}
	n.up = false
	return nil
}
