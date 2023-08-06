package handler

import (
	"context"
	"crypto/hmac"
	"errors"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/usecase/storage"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/gynshu-one/go-metric-collector/proto"
	"github.com/gynshu-one/go-metric-collector/repos/postgres"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"strings"
	"time"
)

type metricServer struct {
	proto.UnimplementedMetricServiceServer
	storage storage.ServerStorage
	dbConn  postgres.DBConn
}

func NewMetricServer(storage storage.ServerStorage, dbConn postgres.DBConn) proto.MetricServiceServer {
	return &metricServer{storage: storage, dbConn: dbConn}
}

func (s *metricServer) Live(ctx context.Context, req *emptypb.Empty) (*proto.LiveResponse, error) {
	resp := &proto.LiveResponse{
		Message: "OK",
	}
	return resp, nil
}

func (s *metricServer) ValueJSON(ctx context.Context, req *proto.ValueRequest) (*proto.MetricResponse, error) {
	var input entity.Metrics
	input.ID = req.GetMetricName()
	input.MType = req.GetMetricType()

	log.Debug().Interface("Request ValueJson Input: %s", input)

	err := getPreCheck(&input)
	if err != nil {
		return nil, handleCustomError(err)
	}
	output := s.storage.Get(input.ID)
	if output == nil {
		return nil, status.Error(codes.NotFound, entity.ErrMetricNotFound.Error())
	}
	output.CalculateHash(config.GetConfig().Key)
	log.Debug().Interface("Request ValueJson Output: %s", output)
	return &proto.MetricResponse{Metric: tools.MarshalMetric(output)}, nil
}

func (s *metricServer) Value(ctx context.Context, req *proto.ValueRequest) (*proto.ValueResponse, error) {
	var input entity.Metrics
	input.ID = req.GetMetricName()
	input.MType = req.GetMetricType()

	log.Debug().Interface("Request Value Input: %s", input)

	err := getPreCheck(&input)
	if err != nil {
		return nil, handleCustomError(err)
	}
	output := s.storage.Get(input.ID)
	if output == nil {
		return nil, status.Error(codes.NotFound, entity.ErrMetricNotFound.Error())
	}
	if output.Value != nil {
		return &proto.ValueResponse{Value: strconv.FormatFloat(
			*output.Value, 'f',
			s.storage.GetFltPrc(input.ID),
			64)}, nil
	} else if output.Delta != nil {
		return &proto.ValueResponse{Value: strconv.FormatInt(*output.Delta, 10)}, nil
	}
	return nil, nil
}

func (s *metricServer) UpdateMetricsJSON(ctx context.Context, req *proto.UpdateMetricsJSONRequest) (*proto.MetricResponse, error) {
	input := tools.UnmarshalMetric(req.GetMetric())

	log.Debug().Interface("Request UpdateMetricsJson Input: %s", input)

	err := setPreCheck(input)
	if err != nil {
		return nil, handleCustomError(err)
	}
	output := s.storage.Set(input)
	if output == nil {
		return nil, status.Error(codes.InvalidArgument, entity.ErrNameTypeMismatch.Error())
	}
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		s.storage.Dump(ctx)
	}
	output.CalculateHash(config.GetConfig().Key)
	return &proto.MetricResponse{Metric: tools.MarshalMetric(output)}, nil
}

func (s *metricServer) UpdateMetric(ctx context.Context, req *proto.UpdateMetricRequest) (*proto.MetricResponse, error) {
	var input entity.Metrics
	input.ID = req.GetMetricName()
	input.MType = req.GetMetricType()
	metricValue := req.GetMetricValue()

	log.Debug().Interface("Request UpdateMetric Input: %s", input)

	switch input.MType {
	case entity.GaugeType:
		val, err_ := strconv.ParseFloat(metricValue, 64)
		if err_ != nil {
			return nil, status.Error(codes.InvalidArgument, entity.ErrNameTypeMismatch.Error())
		}
		input.Value = &val
	case entity.CounterType:
		val, err_ := strconv.ParseInt(metricValue, 10, 64)
		if err_ != nil {
			return nil, status.Error(codes.InvalidArgument, entity.ErrNameTypeMismatch.Error())
		}
		input.Delta = &val
	default:
		input.Delta = nil
	}
	err := setPreCheck(&input)
	if err != nil {
		return nil, handleCustomError(err)
	}
	s.storage.SetFltPrc(input.ID, metricValue)
	output := s.storage.Set(&input)
	if output == nil {
		return nil, status.Error(codes.InvalidArgument, entity.ErrNameTypeMismatch.Error())
	}
	s.storage.SetFltPrc(input.ID, metricValue)
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		s.storage.Dump(ctx)
	}
	output.CalculateHash(config.GetConfig().Key)
	return &proto.MetricResponse{Metric: tools.MarshalMetric(output)}, nil
}

func (s *metricServer) BulkUpdateJSON(ctx context.Context, req *proto.BulkUpdateJSONRequest) (*proto.BulkUpdateResponse, error) {
	var input []*entity.Metrics
	for i := range req.GetMetrics() {
		input = append(input, tools.UnmarshalMetric(req.GetMetrics()[i]))
	}
	log.Debug().Interface("Request BulkUpdateJson %d Inputs", len(input))
	var inputMapper = make(map[string]*entity.Metrics)
	for i := range input {
		err := setPreCheck(input[i])
		if err != nil {
			log.Error().Err(err).Msg("Some of the input metrics are invalid")
			continue
		}
		val := s.storage.Set(input[i])
		if val == nil {
			log.Error().Err(err).Interface("Some of the input metrics are invalid %s", entity.ErrUnableToStore)
			continue
		}
		inputMapper[input[i].ID] = val
	}
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		s.storage.Dump(ctx)
	}

	var output []entity.Metrics
	for i := range inputMapper {
		inputMapper[i].CalculateHash(config.GetConfig().Key)
		output = append(output, *inputMapper[i])
	}
	response := &proto.BulkUpdateResponse{
		Metrics: []*proto.Metric{},
	}
	for i := range output {
		response.Metrics = append(response.Metrics, tools.MarshalMetric(&output[i]))
	}
	return response, nil
}

func (s *metricServer) PingDB(ctx context.Context, req *emptypb.Empty) (*proto.PingDBResponse, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := s.dbConn.Ping(c)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return &proto.PingDBResponse{Message: "Pong"}, nil
}

// getPreCheck checks if the metric is valid for GET request
// returns predefined error if not
func getPreCheck(m *entity.Metrics) error {
	m.MType = strings.ToLower(m.MType)
	if m.ID == "" {
		return entity.ErrMetricNameNotProvided
	}
	if m.MType == "" {
		return entity.ErrMetricTypeNotProvided
	}
	switch m.MType {
	case entity.GaugeType, entity.CounterType:
	default:
		return entity.ErrInvalidType
	}
	return nil
}

func handleCustomError(err error) error {
	switch {
	case errors.Is(err, entity.ErrInvalidType):
		return status.Error(codes.Unimplemented, err.Error())
	case errors.Is(err, entity.ErrTypeValueMismatch), errors.Is(err, entity.ErrInvalidHash):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, entity.ErrDBConnError):
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.NotFound, err.Error())
	}
}

func setPreCheck(m *entity.Metrics) error {
	m.MType = strings.ToLower(m.MType)
	switch m.MType {
	case entity.GaugeType, entity.CounterType:
		if m.MType == entity.GaugeType && m.Value == nil {
			return entity.ErrTypeValueMismatch
		} else if m.MType == entity.CounterType && m.Delta == nil {
			return entity.ErrTypeValueMismatch
		}
	default:
		return entity.ErrInvalidType
	}
	if m.ID == "" {
		return entity.ErrMetricNameNotProvided
	}
	if config.GetConfig().Key != "" {
		inputHash := m.Hash
		m.CalculateHash(config.GetConfig().Key)
		if !hmac.Equal([]byte(inputHash), []byte(m.Hash)) {
			log.Debug().Msgf("Hash mismatch: %s != %s on %s", inputHash, m.Hash, m.String())
			return entity.ErrInvalidHash
		}
	}
	return nil
}
