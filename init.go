package connman

var (
	wifiInterface string = ""
	lastStatus    []byte
)

func init() {
	var err error
	wifiInterface, err = getWifiInterface()
	if err != nil {
		panic(err)
	}
}
