package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandleMetrics(t *testing.T) {
	type args struct {
		method string
		target string
	}
	tests := []struct {
		name               string
		args               args
		expectedStatusCode int
	}{
		{
			name: "Wrong method",
			args: args{
				method: "GET",
				target: "/update/",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Metric type is empty",
			args: args{
				//w 	: httptest.NewRecorder(),
				//r 	: httptest.NewRequest("POST", "/update/gauge/", nil),
				method: "POST",
				target: "/update/",
			},
			expectedStatusCode: http.StatusNotImplemented,
		},
		{
			name: "Metric type is wrong",
			args: args{
				method: "POST",
				target: "/update/invalid",
			},
			expectedStatusCode: http.StatusNotImplemented,
		},
		{
			name: "Metric name is empty",
			args: args{
				method: "POST",
				target: "/update/gauge/",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Metric value is empty",
			args: args{
				method: "POST",
				target: "/update/gauge/cpu",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Invalid metric value",
			args: args{
				method: "POST",
				target: "/update/gauge/cpu/invalid",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Valid metric value",
			args: args{
				method: "POST",
				target: "/update/gauge/cpu/0.5",
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleMetrics)
			r := httptest.NewRequest(tt.args.method, tt.args.target, nil)
			h.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()
			HandleMetrics(w, r)
			if res.StatusCode != tt.expectedStatusCode {
				t.Fatalf("Expected status code %d, got %d", tt.expectedStatusCode, res.StatusCode)

				//t.Errorf("Body: %s", w.Body.String())
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		defaultAddr *string
		defaultPort *string
	}
	addr := "local"
	port := "80"
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "Create a new server",
			args: args{
				defaultAddr: &addr,
				defaultPort: &port,
			},
			want: &Server{
				defaultAddr: &addr,
				defaultPort: &port,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.defaultAddr, tt.args.defaultPort); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestServer_Start(t *testing.T) {
//	type fields struct {
//		defaultAddr *string
//		defaultPort *string
//	}
//	addr := "localhost"
//	port := "8009"
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "Start the server",
//			fields: fields{
//				defaultAddr: &addr,
//				defaultPort: &port,
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Server{
//				defaultAddr: tt.fields.defaultAddr,
//				defaultPort: tt.fields.defaultPort,
//			}
//			go func() {
//				s.Start()
//			}()
//			// make sure the server is running
//			r := httptest.NewRequest("GET", "http://localhost:8009/live", nil)
//			w := httptest.NewRecorder()
//			http.DefaultServeMux.ServeHTTP(w, r)
//			res := w.Result()
//			defer res.Body.Close()
//			if res.StatusCode != http.StatusOK {
//				t.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
//			}
//		})
//	}
//}
