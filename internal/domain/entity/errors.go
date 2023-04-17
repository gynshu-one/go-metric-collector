package entity

import "errors"

var (
	BulkReport            = errors.New("bulk report error")
	UnableToStore         = errors.New("unable to store")
	InvalidType           = errors.New("invalid type")
	TypeValueMismatch     = errors.New("type and value mismatch")
	NameTypeMismatch      = errors.New("name and type you have sent mismatch with the one in the storage")
	MetricTypeNotProvided = errors.New("metric type not provided")
	MetricNameNotProvided = errors.New("metric name not provided")
	MetricNotFound        = errors.New("metric not found")
	InvalidHash           = errors.New("invalid hash")
	HashNotProvided       = errors.New("env var KEY is set but hash is missing")
	DBConnError           = errors.New("db connection error")
	InvalidMetric         = errors.New("invalid metric")
	EmptyMetric           = errors.New("empty metric")
)
