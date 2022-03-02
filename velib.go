package velib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var v *Velib

type Velib struct {
	// States
	ready  bool
	failed bool

	// Configuration
	verbose bool
	port    int

	// Functions
	process func(http.ResponseWriter, *http.Request)
}

func Init() {
	announce("Veritone Engines Library (VELIB) Initialized")
	v = New()
}

func New() *Velib {
	announce("Creating new instance")
	v := new(Velib)
	v.ready = false
	v.failed = false
	v.verbose = false
	v.port = 8080
	v.process = func(responseWriter http.ResponseWriter, request *http.Request) {}
	return v
}

func Run() {
	port := v.getPort()
	l("Starting Engine Server on port", port)
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), server())
	if err != nil {
		l("Error starting server: ")
		fmt.Fprintf(os.Stderr, "exif: %s", err)
		os.Exit(1)
	}
}

func server() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("/ready", handleReady)
	s.HandleFunc("/process", handleProcess)
	return s
}

// the /ready endpoint is used to check if an engine is ready to process
// 200 = Ready to Process
// 503 = Engine is not ready
// 500 = Engine has failed
func handleReady(w http.ResponseWriter, r *http.Request) {
	var readyStatus int
	if v.getReady() {
		readyStatus = http.StatusOK // 200
	} else {
		readyStatus = http.StatusServiceUnavailable // 503
	}

	if v.getFailed() {
		readyStatus = http.StatusInternalServerError // 500
	}

	l("/ready ", readyStatus)
	w.WriteHeader(readyStatus)
}

func handleProcess(responseWriter http.ResponseWriter, request *http.Request) {
	l("/process")
	request.ParseMultipartForm(512 * 1024 * 1024)
	start := time.Now()

	// inject into this part of the process
	v.process(responseWriter, request)

	duration := time.Since(start)
	l("Processing took ", duration)
	heartbeatCallback := request.FormValue("heartbeatWebhook")
	sendHeartbeat(heartbeatCallback, "complete", map[string]string{
		"processingDuration": duration.String(),
	})
}

func sendHeartbeat(callback string, status string, info map[string]string) {
	// heartbeat may not be present in all circumstances, like during testing
	if callback == "" {
		return
	}

	// prepare the body of the heartbeat
	bodyMap := map[string]interface{}{
		"status":  status,
		"infoMsg": info,
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		log.Printf("Unable to marshal the heartbeat body: %s\n", err.Error())
		return
	}

	// post the callback
	resp, err := http.Post(callback, "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		log.Printf("Unable to send heartbeat to aiWARE: %s\n", err.Error())
		return
	}

	// consume the response
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
}

func announce(text string) {
	log.Println("[VELIB] ", text)
}

func l(text ...interface{}) {
	if v.getVerbose() {
		log.Println("[VELIB] ", text)
	}
}
