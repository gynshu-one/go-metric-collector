package handlers

import (
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"log"
	"net/http"
	"runtime"
	"time"
)

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	serverAddr     string
	metrics        *storage.MemStorage
}

func NewAgent(pollInterval, reportInterval time.Duration, serverAddr string) *Agent {
	return &Agent{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddr:     serverAddr,
		metrics:        storage.InitStorage(),
	}
}

// Poll polls runtime metrics and reports them to the server by calling Report()
func (a *Agent) Poll() {
	memStats := &runtime.MemStats{}
	a.metrics.Counter["PollCount"] = 0
	// ReadRuntime runtime metrics
	go func() {
		for {
			runtime.ReadMemStats(memStats)
			a.metrics.ReadRuntime(memStats)
			// Sleep for poll interval
			time.Sleep(a.pollInterval)
		}
	}()
	// Report
	for {
		time.Sleep(a.reportInterval)
		a.Report()
	}

}

func (a *Agent) Report() {
	a.metrics.ApplyToAll(a.MakeReport)
	a.metrics.Counter["PollCount"] = 0
}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func (a *Agent) MakeReport(t, n, v string) {
	// Create a new request
	req, err := http.NewRequest("POST",
		a.serverAddr+"/update/"+t+"/"+n+"/"+v, nil)
	if err != nil {
		log.Fatalf("\nUpdate request maker failed with error: %e\n", err)

	}
	// Set the content type
	req.Header.Set("Content-Type", "text/plain")
	// Create a client
	client := &http.Client{}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("\nUpdate request failed with error: %e\n", err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("\nUpdate request failed with error: %e\n", err)

	}
}
