package service

import (
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"math/rand"
	"testing"
)

func populate(numberOfElements int, service *memService) {
	metrics := make([]*entity.Metrics, 0, numberOfElements)
	for i := 0; i < numberOfElements; i++ {
		randN := rand.Intn(999999999)
		randN64 := int64(randN)
		metrics = append(metrics, entity.NewMetrics("test"+fmt.Sprintf("_%d", i), entity.CounterType, randN64))
	}
	for _, m := range metrics {
		service.Set(m)
	}
	return
}

func BenchmarkStorageService_Set(b *testing.B) {
	service := NewMemService()
	metrics := make([]*entity.Metrics, 0, b.N)
	for i := 0; i < b.N; i++ {
		randN := rand.Intn(999999999)
		randN64 := int64(randN)
		metrics = append(metrics, entity.NewMetrics("test"+fmt.Sprintf("_%d", randN), entity.CounterType, randN64))
	}
	k := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Set(metrics[k])
		k++
	}

}

func BenchmarkStorageService_Get(b *testing.B) {
	service := NewMemService()
	populate(b.N, service)
	k := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		f := entity.Metrics{
			ID: "test" + fmt.Sprintf("_%d", k),
		}
		b.StartTimer()
		service.Get(f.ID)
		k++
	}
}

func BenchmarkStorageService_ApplyToAll(b *testing.B) {
	service := NewMemService()
	populate(10000, service)
	k := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ApplyToAll(func(m *entity.Metrics) {
			m.String()
		})
		//service.GetAll()
		k++
	}

}

func BenchmarkStorageService_GetAll(b *testing.B) {
	service := NewMemService()
	populate(10000, service)
	k := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAll()
		k++
	}

}
