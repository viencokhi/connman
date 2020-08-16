package connman

import (
	"fmt"
	"time"
)

var (
	hotspotTime   time.Duration = 1 * time.Minute
	nodeResetTime time.Duration = 1 * time.Minute
	mainSpot      *Spot         = NewSpot("node1234", "")
)

func startUp() {
	lastStatus = status()
	disconnectFromAllConnection()
	// resetNetwork()
	hotspotOn(mainSpot)
	go startToListen()
	time.Sleep(hotspotTime)
	hotspotOff(mainSpot)
	server.Close()
	connectAvailable()
}

func nodeLoop() {
	fmt.Println("node started...")
	time.Sleep(hotspotTime)
}

// Routine automation routine
func Routine() {
	startUp()
	for true {
		nodeLoop()
		time.Sleep(nodeResetTime)
	}
}
