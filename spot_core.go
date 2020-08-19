package connman

import (
	"fmt"
	"strings"
)

func addSpot(name, pass string) error {
	var cmd string
	cmd = fmt.Sprintf(`nmcli con add type wifi ifname "%v" con-name "%v" autoconnect yes ssid "%v"`, wifiInterface, name, name)
	out, err := exe(cmd, "add spot")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}

	cmd = fmt.Sprintf(`nmcli con modify "%v" 802-11-wireless.mode ap 802-11-wireless.band bg ipv4.method shared`, name)
	out, err = exe(cmd, "modify conn spot")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}

	if pass != "" {
		cmd = fmt.Sprintf(`nmcli con modify "%v" wifi-sec.key-mgmt wpa-psk`, name)
		out, err = exe(cmd, "modify pass type spot")
		if err != nil {
			return fmt.Errorf("error:%v, out:%v", err.Error(), out)
		}

		cmd = fmt.Sprintf(`nmcli con modify "%v" wifi-sec.psk "%v"`, name, pass)
		out, err = exe(cmd, "modify pass spot")
		if err != nil {
			return fmt.Errorf("error:%v, out:%v", err.Error(), out)
		}
	}
	return nil
}

// up, down, delete
func setSpot(mode, name string) error {
	cmd := fmt.Sprintf(`nmcli con "%v" "%v"`, mode, name)
	out, err := exe(cmd, fmt.Sprintf(`modify "%v" spot`, mode))
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	return nil
}

//SpotExists spot exists return bool
func spotExists(name string) (bool, error) {
	cmd := fmt.Sprintf(`nmcli -f NAME con sh`)
	out, err := exe(cmd, "spot scan")
	if err != nil {
		return false, fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	tokens := strings.Split(out, "\n")
	// remove column name
	tokens = tokens[1:]
	for _, t := range tokens {
		tc := strings.TrimSpace(t)
		if tc == name {
			return true, nil
		}
	}
	return false, nil
}
