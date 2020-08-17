package connman

import (
	"time"
)

var (
	mainSpotName  string        = "solar"
	signInTime    time.Duration = 1 * time.Minute
	configTime    time.Duration = 5 * time.Minute
	mainSpot      *Spot         = NewSpot(mainSpotName, "")
	closeChannel  chan bool     = make(chan bool, 1)
	signInChannel chan bool     = make(chan bool, 1)
)

//StartUp main startup function.
func StartUp() {
	firstAttemptFailed := false
	attemptCount := 0
	hasInternetConnection := false

	for !hasInternetConnection {
		disconnectFromAllConnection()
		connectAvailable()
		lastStatus = status()
		disconnectFromAllConnection()
		hotspotOn(mainSpot)
		go startToListen()
		if firstAttemptFailed {
			<-closeChannel
		} else {
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
		}
		hotspotOff(mainSpot)
		server.Close()
		connectAvailable()
		hasInternetConnection = connectedToInternet()
		if attemptCount == 0 && !hasInternetConnection {
			firstAttemptFailed = true
		}
		attemptCount++
	}
}
