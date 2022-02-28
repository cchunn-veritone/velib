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
	"time"
)

var v *Velib

type Velib struct {
	ready   bool
	failed  bool
	verbose bool
}

func Init() *Velib {
	v := new(Velib)
	v.ready = false
	v.failed = false
	v.verbose = false

	return v
}

func Run() {
	l("Veritone Engines Library (VELIB)")
	if err := http.ListenAndServe("0.0.0.0:8080", server()); err != nil {
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
	if v.ready {
		readyStatus = http.StatusOK // 200
	} else {
		readyStatus = http.StatusServiceUnavailable // 503
	}

	if v.failed {
		readyStatus = http.StatusInternalServerError // 500
	}

	l("/ready ", readyStatus)
	w.WriteHeader(readyStatus)
}

func handleProcess(responseWriter http.ResponseWriter, request *http.Request) {
	l("/process")
	request.ParseMultipartForm(512 * 1024 * 1024)
	start := time.Now()

	// How do we inject code here to handle the request?

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
