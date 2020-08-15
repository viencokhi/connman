package connman

var wifiInterface string = ""

func init() {
	var err error
	wifiInterface, err = getWifiInterface()
	if err != nil {
		panic(err)
	}
}
