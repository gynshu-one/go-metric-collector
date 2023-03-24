package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLive(t *testing.T) {
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
			name: "Server is live",
			args: args{
				method: "GET",
				target: "/live",
			},
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.GET("/live", Live)
			w := SetUpRouter(r, tt.args.method, tt.args.target)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Live() = %v, want %v", w.Code, tt.expectedStatusCode)

			}
		})
	}
}

func TestUpdateMetrics(t *testing.T) {
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
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Metric type is empty",
			args: args{
				//w 	: httptest.NewRecorder(),
				//r 	: httptest.NewRequest("POST", "/update/gauge/", nil),
				method: "POST",
				target: "/update/",
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Metric type is wrong",
			args: args{
				method: "POST",
				target: "/update/invalid/name/val",
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
			name: "Metric name is empty",
			args: args{
				method: "POST",
				target: "/update/counter/",
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
			r := gin.Default()
			r.POST("/update/:metric_type/:metric_name/:metric_value", UpdateMetrics)
			w := SetUpRouter(r, tt.args.method, tt.args.target)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("UpdateMetrics() = %v, want %v", w.Code, tt.expectedStatusCode)
			}
		})
	}
}
func SetUpRouter(r *gin.Engine, method string, target string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = true
	r.ServeHTTP(w, req)
	return w
}
