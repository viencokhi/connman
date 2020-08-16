package connman

import (
	"fmt"
	"time"
)

var (
	mainSpotName  string        = "solar"
	signInTime    time.Duration = 1 * time.Minute
	configTime    time.Duration = 5 * time.Minute
	nodeResetTime time.Duration = 1 * time.Minute
	mainSpot      *Spot         = NewSpot(mainSpotName, "")
	closeChannel  chan bool     = make(chan bool, 1)
	signInChannel chan bool     = make(chan bool, 1)
)

func startUp() {
	lastStatus = status()
	disconnectFromAllConnection()
	hotspotOn(mainSpot)
	go startToListen()

	select {
	case <-signInChannel:
		select {
		case <-closeChannel:
			break
		case <-time.After(configTime):
			break
		}
		break
	case <-closeChannel:
		break
	case <-time.After(signInTime):
		break
	}

	hotspotOff(mainSpot)
	server.Close()
	connectAvailable()
}

func nodeLoop() {
	fmt.Println("node started...")
	time.Sleep(nodeResetTime)
}

// Routine automation routine
func Routine() {
	startUp()
	for true {
		nodeLoop()
	}
}
