package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAgent_MakeReport(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "test1",
			args: []string{"metric1", "value1"},
			want: "/update/test1/metric1/value1",
		},
		{
			name: "test2",
			args: []string{"metric2", "value2"},
			want: "/update/test2/metric2/value2",
		},
	}

	pollInterval := 2 * time.Millisecond
	reportInterval := 10 * time.Second
	serverAddr := "http://localhost:8080"
	equal := func(t *testing.T, a, b interface{}) {
		t.Helper()
		if !reflect.DeepEqual(a, b) {
			t.Errorf("request path to be %v but got %v", a, b)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				equal(t, req.URL.String(), tt.want)
				rw.WriteHeader(http.StatusOK)
			}))
			serverAddr = server.URL
			defer server.Close()
			agent := NewAgent(pollInterval, reportInterval, serverAddr)
			agent.MakeReport(tt.name, tt.args[0], tt.args[1])
		})
	}
}

//func TestAgent_Report(t *testing.T) {
//	pollInterval := 100 * time.Millisecond
//	reportInterval := 1 * time.Second
//	serverAddr := "localhost:8080"
//	agent := NewAgent(&pollInterval, &reportInterval, &serverAddr)
//
//	agent.Report()
//
//	if agent.metrics.Custom.PollCount != 0 {
//		t.Errorf("Expected poll count to be 0 but got %d", agent.metrics.Custom.PollCount)
//	}
//}

//func TestAgent_Poll(t *testing.T) {
//	pollInterval := 100 * time.Millisecond
//	reportInterval := 1 * time.Second
//	serverAddr := "localhost:8080"
//	agent := NewAgent(&pollInterval, &reportInterval, &serverAddr)
//
//	go agent.Poll()
//	time.Sleep(2 * time.Second)
//
//	if agent.metrics.Custom.PollCount == 0 {
//		t.Errorf("Expected poll count to be non-zero but got 0")
//	}
//}

func TestNewAgent(t *testing.T) {
	pollInterval := time.Minute
	reportInterval := time.Hour
	serverAddr := "http://localhost:8080"
	agent := NewAgent(pollInterval, reportInterval, serverAddr)

	if agent.pollInterval != pollInterval {
		t.Errorf("Expected poll interval to be %v but got %v", pollInterval, agent.pollInterval)
	}
	if agent.reportInterval != reportInterval {
		t.Errorf("Expected report interval to be %v but got %v", reportInterval, agent.reportInterval)
	}
	if agent.serverAddr != serverAddr {
		t.Errorf("Expected server address to be %v but got %v", serverAddr, agent.serverAddr)
	}

	if agent.metrics.Gauge == nil {
		t.Errorf("Expected basic metrics storage to be non-nil but got nil")
	}
	if agent.metrics.Counter == nil {
		t.Errorf("Expected custom metrics storage to be non-nil but got nil")
	}
}
