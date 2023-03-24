package handlers

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	defaultAddr *string
	defaultPort *string
}

func NewServer(defaultAddr *string, defaultPort *string) *Server {
	return &Server{
		defaultAddr: defaultAddr,
		defaultPort: defaultPort,
	}
}

func (s *Server) Start() {
	// Start basic server
	http.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Server is live"))
	})
	http.HandleFunc("/update/", HandleMetrics)
	fmt.Printf("Server started at %s:%s", *s.defaultAddr, *s.defaultPort)
	server := &http.Server{
		Addr: "localhost:8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	// The first and second segment will be " " and "update", so start from the second one
	path := ""
	if strings.HasSuffix(r.URL.Path, "/") {
		path = r.URL.Path[1 : len(r.URL.Path)-1]
	} else {
		path = r.URL.Path[1:]
	}
	segments := strings.Split(path, "/")
	segments = segments[1:]
	fmt.Printf("Segments: %v", segments)
	// check if request is post
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		color.Red("Only POST requests are allowed")
		w.Write([]byte("Only POST requests are allowed"))
		return
	}
	// Check content type
	//if r.Header.Get("Content-Type") != "text/plain" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Header().Set("Content-Type", "application/json")
	//	color.Red("Content-Type must be text/plain")
	//	w.Write([]byte("Content-Type must be text/plain"))
	//	return
	//}

	if len(segments) < 1 || segments[0] == "" {
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Content-Type", "application/json")
		color.Red("Please provide metric type eg. /update/gauge")
		w.Write([]byte("Please provide metric type eg. /update/gauge"))
		return
	}
	// Lowercase the metric type to avoid case sensitivity
	segments[0] = strings.ToLower(segments[0])
	// Check if metric type is valid
	if segments[0] != "gauge" && segments[0] != "counter" {
		w.WriteHeader(http.StatusNotImplemented)
		w.Header().Set("Content-Type", "application/json")
		color.Red("Invalid metric type")
		w.Write([]byte("Invalid metric type"))
		return
	}

	metricType := segments[0]
	if len(segments) < 2 || segments[1] == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		color.Red("\nPlease provide metric name separated by  / eg. /update/gauge/cpu/0.5")
		w.Write([]byte("Please provide metric name  separated by / eg. /update/gauge/cpu/0.5"))
		return
	}
	metricName := segments[1]
	if len(segments) < 3 || segments[2] == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		color.Red("\nPlease provide metric value separated by / eg. /update/gauge/cpu/0.5")
		w.Write([]byte("Please provide metric value separated by / eg. /update/gauge/cpu/0.5"))
		return
	}
	metricValue := segments[2]
	//if !storage.Stor.CheckIfNameExists(metricName) {
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Header().Set("Content-Type", "application/json")
	//	color.Red("\nMetric name is not acceptable")
	//	w.Write([]byte("Metric name is not acceptable"))
	//	return
	//}
	// ReadRuntime
	err := storage.Stor.AddMetric(metricType, metricName, metricValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		color.Red("\nInvalid metric value %s Should be number", metricValue)
		w.Write([]byte("Invalid metric value. Should be number"))
		return
	}
	// Send status code 200
	//fmt.Printf("\nReceived and collected metric: %s %s %s", metricType, metricName, metricValue)
	w.WriteHeader(http.StatusOK)
	//storage.Stor.PrintAll()
	w.Write([]byte(metricType + " " + metricName + " " + metricValue + "\n"))
}
