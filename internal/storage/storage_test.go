package storage

import (
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"testing"
)

func MakeStorage() ServerInterface {
	configs.CFG = &configs.Config{
		Address:       "localhost:8080",
		StoreInterval: 10,
		StoreFile:     "/tmp/devops-metrics-db.json",
		Restore:       false,
	}
	dir := tools.GetProjectRoot()
	// Make temp files dir absolute
	configs.CFG.StoreFile = dir + configs.CFG.StoreFile
	configs.CFG.InitFiles()
	return InitServerStorage()
}
func TestValidateTypeAndValue(t *testing.T) {
	storage := MakeStorage()
	gaugeMetric := Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(1.0),
	}
	assert.True(t, storage.ValidateTypeAndValue(gaugeMetric), "Gauge metric with value should be valid")

	counterMetric := Metrics{
		ID:    "TestCounter",
		MType: "counter",
		Delta: tools.Int64Ptr(1),
	}
	assert.True(t, storage.ValidateTypeAndValue(counterMetric), "Counter metric with delta should be valid")

	invalidGaugeMetric := Metrics{
		ID:    "InvalidGauge",
		MType: "gauge",
		Delta: tools.Int64Ptr(1),
	}
	assert.False(t, storage.ValidateTypeAndValue(invalidGaugeMetric), "Gauge metric with delta should be invalid")

	invalidCounterMetric := Metrics{
		ID:    "InvalidCounter",
		MType: "counter",
		Value: tools.Float64Ptr(1.0),
	}
	assert.False(t, storage.ValidateTypeAndValue(invalidCounterMetric), "Counter metric with value should be invalid")
}

func TestUpdateMetric(t *testing.T) {
	storage := MakeStorage()
	gaugeMetric := Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(1.0),
	}
	updatedGauge := storage.UpdateMetric(gaugeMetric)
	assert.Equal(t, gaugeMetric, updatedGauge, "Gauge metric should be updated correctly")

	counterMetric := Metrics{
		ID:    "TestCounter",
		MType: "counter",
		Delta: tools.Int64Ptr(1),
	}
	updatedCounter := storage.UpdateMetric(counterMetric)
	assert.Equal(t, counterMetric, updatedCounter, "Counter metric should be updated correctly")

	newCounterMetric := Metrics{
		ID:    "TestCounter",
		MType: "counter",
		Delta: tools.Int64Ptr(2),
	}
	expectedUpdatedCounter := Metrics{
		ID:    "TestCounter",
		MType: "counter",
		Delta: tools.Int64Ptr(3),
	}
	updatedCounter = storage.UpdateMetric(newCounterMetric)
	assert.Equal(t, expectedUpdatedCounter, updatedCounter, "Existing counter metric should be updated correctly")
}

func TestCheckMetricType(t *testing.T) {
	storage := MakeStorage()
	assert.True(t, storage.CheckMetricType("gauge"), "Gauge should be a valid metric type")
	assert.True(t, storage.CheckMetricType("counter"), "Counter should be a valid metric type")
	assert.False(t, storage.CheckMetricType("invalidType"), "Invalid type should be invalid")
}
func TestFindMetricByName(t *testing.T) {
	storage := MakeStorage()
	// Test not found metric
	_, ok := storage.FindMetricByName("NonExistent")
	assert.False(t, ok, "Non-existent metric should not be found")

	// Test found metric
	gaugeMetric := Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(1.0),
	}
	storage.UpdateMetric(gaugeMetric)
	foundMetric, ok := storage.FindMetricByName("TestGauge")
	assert.True(t, ok, "Existing metric should be found")
	assert.Equal(t, gaugeMetric, foundMetric, "Found metric should match the expected metric")
}

func TestRandomValue(t *testing.T) {
	storage := InitAgentStorage()
	storage.RandomValue()

	metric, ok := storage.FindMetricByName("RandomValue")
	assert.True(t, ok, "RandomValue metric should exist after calling RandomValue")
	assert.NotNil(t, metric.Value, "RandomValue metric should have a value")
}

func TestAddPollCount(t *testing.T) {
	storage := InitAgentStorage()
	storage.AddPollCount()

	metric, ok := storage.FindMetricByName("PollCount")
	assert.True(t, ok, "PollCount metric should exist after calling AddPollCount")
	assert.NotNil(t, metric.Delta, "PollCount metric should have a delta")
	assert.Equal(t, int64(1), *metric.Delta, "Initial PollCount metric delta should be 1")

	storage.AddPollCount()
	metric, _ = storage.FindMetricByName("PollCount")
	assert.Equal(t, int64(2), *metric.Delta, "PollCount metric delta should be incremented to 2")
}

func TestCheckIfNameExists(t *testing.T) {
	storage := MakeStorage()
	assert.False(t, storage.CheckIfNameExists("NonExistent"), "Non-existent metric should return false")

	gaugeMetric := Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(1.0),
	}
	storage.UpdateMetric(gaugeMetric)
	assert.True(t, storage.CheckIfNameExists("TestGauge"), "Existing metric should return true")
}

func TestReadRuntime(t *testing.T) {
	storage := InitAgentStorage()
	storage.ReadRuntime()

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	input := reflect.ValueOf(memStats).Elem()
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	for i := 0; i < input.NumField(); i++ {
		fieldName := input.Type().Field(i).Name
		if !tools.Contains(defaultExclusion, fieldName) {
			if _, ok := storage.FindMetricByName(fieldName); !ok {
				t.Errorf("Expected runtime metric %s not found in MemStorage", fieldName)
			}
		}
	}
}
