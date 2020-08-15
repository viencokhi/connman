package connman

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ecoshub/jin"
	"github.com/gorilla/mux"
)

const (
	masterPort string = "80"
)

var (
	server     *http.Server
	mainRouter *mux.Router
)

func routeCreation() {
	mainRouter = mux.NewRouter()
	mainRouter.HandleFunc("/", rootHandle)
	mainRouter.HandleFunc("/state", stateHandle)
	mainRouter.HandleFunc("/add", addHandle)
	mainRouter.HandleFunc("/remove", addHandle)
}

func serverCreation() {
	server = &http.Server{
		Handler: mainRouter,
		Addr:    "localhost:" + masterPort,
	}
}

//StartToListen main listen function
func startToListen() {
	routeCreation()
	serverCreation()
	server.ListenAndServe()
	fmt.Println("un expected shutdown")
}

func headerMiddleware(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, cache-control")
	w.Header().Set("Content-Type", "application/json")
	return w
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message":"welcome to solarnode"}`))
}

func stateHandle(w http.ResponseWriter, r *http.Request) {
	w = headerMiddleware(w)
	w.Write(lastStatus)
}

func addHandle(w http.ResponseWriter, r *http.Request) {
	w = headerMiddleware(w)
	json, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(`{"error":"JSON parse error", "success":false}`))
		return
	}
	network, err := jin.GetString(json, "network")
	if err != nil {
		w.Write([]byte(`{"error":"'network' field not found", "success":false}`))
		return
	}
	pass, err := jin.GetString(json, "password")
	if err != nil {
		w.Write([]byte(`{"error":"'password' field not found", "success":false}`))
		return
	}
	err = saveNetwork(network, pass)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%v", "success":false}`, err.Error())))
		return
	}
	w.Write([]byte(`{"success":true, "error":null}`))
}

func removeHandle(w http.ResponseWriter, r *http.Request) {
	w = headerMiddleware(w)
	json, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(`{"error":"JSON parse error", "success":false}`))
		return
	}
	network, err := jin.GetString(json, "network")
	if err != nil {
		w.Write([]byte(`{"error":"'network' field not found", "success":false}`))
		return
	}
	err = removeNetwork(network)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":"%v", "success":false}`, err.Error())))
		return
	}
	w.Write([]byte(`{"success":true, "error":null}`))
}
