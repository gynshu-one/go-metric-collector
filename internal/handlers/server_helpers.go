package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"net/http"
	"strings"
)

func getPreCheck(m *storage.Metrics) error {
	m.MType = strings.ToLower(m.MType)
	if m.ID == "" {
		return errors.New(storage.MetricNameNotProvided)
	}
	if m.MType == "" {
		return errors.New(storage.MetricTypeNotProvided)
	}
	switch m.MType {
	case storage.GaugeType, storage.CounterType:
	default:
		return errors.New(storage.InvalidType)
	}
	return nil
}
func setPreCheck(m *storage.Metrics) error {
	m.MType = strings.ToLower(m.MType)
	switch m.MType {
	case storage.GaugeType, storage.CounterType:
		if m.MType == storage.GaugeType && m.Value == nil {
			return errors.New(storage.TypeValueMismatch)
		}
		if m.MType == storage.CounterType && m.Delta == nil {
			return errors.New(storage.TypeValueMismatch)
		}
	default:
		return errors.New(storage.InvalidType)
	}
	if m.ID == "" {
		return errors.New(storage.MetricNameNotProvided)
	}

	return nil
}
func handleCustomError(ctx *gin.Context, err error) {
	switch err.Error() {
	case storage.InvalidType:
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
		return
	case storage.TypeValueMismatch:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	default:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
}
func generateHTMLTable(M storage.MemActions) []string {
	var table []string
	M.ApplyToAll(func(m *storage.Metrics) {
		val := ""
		if m.Value != nil {
			val = fmt.Sprintf("%f", *m.Value)
		}
		if m.Delta != nil {
			val = fmt.Sprintf("%d", *m.Delta)
		}
		table = append(table, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>",
			m.MType, m.ID, val))
	})
	return table
}
