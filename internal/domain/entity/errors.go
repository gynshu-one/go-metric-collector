package entity

import "errors"

var (
	ErrBulkReport            = errors.New("bulk report error")
	ErrUnableToStore         = errors.New("unable to store")
	ErrInvalidType           = errors.New("invalid type")
	ErrTypeValueMismatch     = errors.New("type and value mismatch")
	ErrNameTypeMismatch      = errors.New("name and type you have sent mismatch with the one in the storage")
	ErrMetricTypeNotProvided = errors.New("metric type not provided")
	ErrMetricNameNotProvided = errors.New("metric name not provided")
	ErrMetricNotFound        = errors.New("metric not found")
	ErrInvalidHash           = errors.New("invalid hash")
	ErrDBConnError           = errors.New("db connection error")
	ErrInvalidMetric         = errors.New("invalid metric")
)
