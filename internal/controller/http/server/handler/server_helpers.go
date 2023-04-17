package handler

import (
	"crypto/hmac"
	"fmt"
	"github.com/gin-gonic/gin"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/usecase/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func getPreCheck(m *entity.Metrics) error {
	m.MType = strings.ToLower(m.MType)
	if m.ID == "" {
		return entity.MetricNameNotProvided
	}
	if m.MType == "" {
		return entity.MetricTypeNotProvided
	}
	switch m.MType {
	case entity.GaugeType, entity.CounterType:
	default:
		return entity.InvalidType
	}
	return nil
}
func setPreCheck(m *entity.Metrics) error {
	m.MType = strings.ToLower(m.MType)
	switch m.MType {
	case entity.GaugeType, entity.CounterType:
		if m.MType == entity.GaugeType && m.Value == nil {
			return entity.TypeValueMismatch
		} else if m.MType == entity.CounterType && m.Delta == nil {
			return entity.TypeValueMismatch
		}
	default:
		return entity.InvalidType
	}
	if m.ID == "" {
		return entity.MetricNameNotProvided
	}
	if config.GetConfig().Key != "" {
		inputHash := m.Hash
		m.CalculateHash(config.GetConfig().Key)
		if !hmac.Equal([]byte(inputHash), []byte(m.Hash)) {
			log.Debug().Msgf("Hash mismatch: %s != %s on %s", inputHash, m.Hash, m.String())
			return entity.InvalidHash
		}
	}
	return nil
}
func handleCustomError(ctx *gin.Context, err error) {
	switch err {
	case entity.InvalidType:
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
		return
	case entity.TypeValueMismatch, entity.InvalidHash:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case entity.DBConnError:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	default:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
}
func generateHTMLTable(M storage.ServerStorage) []string {
	var table []string
	M.ApplyToAll(func(m *entity.Metrics) {
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
