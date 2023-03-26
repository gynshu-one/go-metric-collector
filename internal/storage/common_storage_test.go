package storage

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMemStorage_AddMetric(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		tp    string
		name  string
		value string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantGauge   map[string]float64
		wantCounter map[string]int64
	}{
		{
			name: "AddValidMetric",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
			args: args{
				tp:    "Gauge",
				name:  "Alloc",
				value: "100",
			},
			wantErr: false,
			wantGauge: map[string]float64{
				"Alloc": 100,
			},
		},
		{
			name: "AddFalseMetric",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
			args: args{
				tp:    "Gauge",
				name:  "Alloc",
				value: "all",
			},
			wantErr: true,
		},
		{
			name: "AddValidCounter",
			fields: fields{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
			args: args{
				tp:    "Counter",
				name:  "PollCount",
				value: "100",
			},
			wantErr: false,
			wantCounter: map[string]int64{
				"PollCount": 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			M := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			if err := M.AddMetric(tt.args.tp, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("AddMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantGauge != nil {
				for k, v := range tt.wantGauge {
					if M.Gauge[k] != v {
						t.Errorf("AddMetric() error expected %v, got %v", v, M.Gauge[k])
					}
				}
			}
			if tt.wantCounter != nil {
				for k, v := range tt.wantCounter {
					if M.Counter[k] != v {
						t.Errorf("AddMetric() error expected %v, got %v", v, M.Counter[k])
					}
				}
			}
		})
	}
}

func TestMemStorage_ApplyToAll(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		f       ApplyToAll
		exclude []string
	}
	var Result string
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ApplyToAll",
			fields: fields{
				Gauge: map[string]float64{
					"Alloc": 100,
				},
				Counter: map[string]int64{
					"PollCount": 100,
				},
			},
			args: args{
				f: func(tp string, name string, value string) {
					Result += fmt.Sprintf("type: %s, name: %s, value: %s", tp, name, value)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			M := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			M.ApplyToAll(tt.args.f, tt.args.exclude...)
			if Result != "type: Gauge, name: Alloc, value: 100type: Counter, name: PollCount, value: 100" {
				t.Errorf("ApplyToAll() error expected %v, got %v", "type: Gauge, name: Alloc, value: 100type: Counter, name: PollCount, value: 100", Result)
			}
		})
	}
}

func TestMemStorage_CheckIfNameExists(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "CheckIfNameExists",
			fields: fields{
				Gauge: map[string]float64{
					"Alloc": 100,
				},
			},
			args: args{
				name: "Alloc",
			},
			want: true,
		},
		{
			name: "CheckIfNameDoesNotExists",
			fields: fields{
				Gauge: map[string]float64{
					"Alloc": 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			M := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			if got := M.CheckIfNameExists(tt.args.name); got != tt.want {
				t.Errorf("CheckIfNameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestMemStorage_PrintAll(t *testing.T) {
//	type fields struct {
//		Gauge   map[string]float64
//		Counter map[string]int64
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			M := &MemStorage{
//				Gauge:   tt.fields.Gauge,
//				Counter: tt.fields.Counter,
//			}
//			M.PrintAll()
//		})
//	}
//}

func TestMemStorage_ReadRuntime(t *testing.T) {
	type fields struct {
		Gauge   map[string]float64
		Counter map[string]int64
	}
	type args struct {
		memStats *runtime.MemStats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ReadRuntime",
			fields: fields{
				Gauge: map[string]float64{
					"Alloc":     100,
					"HeapAlloc": 100,
				},
				Counter: map[string]int64{
					"PollCount": 100,
				},
			},
			args: args{
				memStats: &runtime.MemStats{
					Alloc:     200,
					HeapAlloc: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			M := &MemStorage{
				Gauge:   tt.fields.Gauge,
				Counter: tt.fields.Counter,
			}
			M.ReadRuntime(tt.args.memStats)
			if M.Gauge["Alloc"] != 200 {
				t.Errorf("ReadRuntime() error expected %v, got %v", 200, M.Gauge["Alloc"])
			}
			if _, ok := M.Gauge["HeapAlloc"]; !ok {
				t.Fatalf("ReadRuntime() error expected %v, got %v", true, ok)
			}
			if M.Gauge["HeapAlloc"] != 1 {
				t.Errorf("ReadRuntime() error expected %v, got %v", 1, M.Gauge["HeapAlloc"])
			}
		})
	}
}
