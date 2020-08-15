package connman

import (
	"fmt"
	"time"
)

var (
	hotspotTime   time.Duration = 5 * time.Minute
	nodeResetTime time.Duration = 1 * time.Minute
	mainSpot      *Spot         = NewSpot("solarnode", "")
)

func startUp() {
	resetNetwork()
	lastStatus = status()
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
