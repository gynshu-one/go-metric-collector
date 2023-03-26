package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"log"
	"net/http"
	"time"
)

type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerAddr     string
	Metrics        storage.MemInterface
}

func NewAgent(pollInterval, reportInterval time.Duration, serverAddr string) *Agent {
	return &Agent{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		ServerAddr:     serverAddr,
		Metrics:        storage.InitStorage(),
	}
}

// Poll polls runtime Metrics and reports them to the server by calling Report()
func (a *Agent) Poll() {
	// ReadRuntime runtime Metrics
	go func() {
		for {
			a.Metrics.AddPollCount()
			a.Metrics.ReadRuntime()
			// Sleep for poll interval
			time.Sleep(a.PollInterval)
		}
	}()
	// Report
	for {
		time.Sleep(a.ReportInterval)
		a.Metrics.RandomValue()
		a.Report()
	}

}

func (a *Agent) Report() {
	// check if the metric is presented in MemStorage
	a.Metrics.ApplyToAll(a.MakeReport)
	a.Metrics.PrintAll()
}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func (a *Agent) MakeReport(m storage.Metrics) {
	// Create a new request
	bd, err := json.Marshal(&m)
	if err != nil {
		log.Fatalf("\nFailed to marshal Metrics: %e\n", err)
	}
	newReader := bytes.NewReader(bd)
	req, err := http.NewRequest("POST",
		a.ServerAddr+"/update/", newReader)
	if err != nil {
		log.Fatalf("\nUpdate request maker failed with error: %e\n", err)

	}
	// Set the content type
	req.Header.Set("Content-Type", "application/json")
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
