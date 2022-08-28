package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	cm "gitlab.com/cronolabs/cable/internal/message"
	cp "gitlab.com/cronolabs/cable/internal/proxy"
	ct "gitlab.com/cronolabs/cable/internal/types"
	ws "gitlab.com/cronolabs/cable/internal/websocket"
)

var HealthResponseBody ct.Json = ct.Json{"status": "UP"}
var BadBodyResponseBody ct.Json = ct.Json{"error": "Invalid Data"}

func RunServer() {
	r := router()
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/sip", putSip).
		Methods("PUT").
		Headers("Content-Type", "application/json")
	router.HandleFunc("/health", getHealth).
		Methods("GET")
	router.HandleFunc("/sub", ws.Subscribe)
	return router
}

func getHealth(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	jsonBody, _ := json.Marshal(HealthResponseBody)
	fmt.Fprint(w, string(jsonBody))
}

func putSip(w http.ResponseWriter, req *http.Request) {
	var jsonData ct.Json

	if err := json.NewDecoder(req.Body).Decode(&jsonData); err != nil {
		jsonBody, _ := json.Marshal(BadBodyResponseBody)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(jsonBody))
		return
	}

	jsonBody, _ := json.Marshal(jsonData)

	msg := cm.NewMessage(&jsonData)
	go cp.Relay(msg)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, string(jsonBody))
}
