package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"log"
	"time"
)

var Client = resty.New()

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
		go func() {
			a.Metrics.RandomValue()
			a.Report()
		}()
	}

}

func (a *Agent) Report() {
	// check if the metric is presented in MemStorage

	a.Metrics.ApplyToAll(a.MakeReport)
	//a.Metrics.PrintAll()
}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func (a *Agent) MakeReport(m storage.Metrics) {
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(a.ServerAddr + "/update/")

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Response: %v", resp)
}
