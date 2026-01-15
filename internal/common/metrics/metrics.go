// Package metrics provides centralized metrics collection
package metrics

import (
	"context"
	"time"
)

// Collector interface defines metrics collection operations
type Collector interface {
	RecordLatency(ctx context.Context, metricName string, duration time.Duration)
	RecordFailure(ctx context.Context, metricName string)
	RecordSuccess(ctx context.Context, metricName string)
}

// DefaultCollector implements a no-op metrics collector
type DefaultCollector struct{}

func (d *DefaultCollector) RecordLatency(ctx context.Context, metricName string, duration time.Duration) {
	// No-op implementation
}

func (d *DefaultCollector) RecordFailure(ctx context.Context, metricName string) {
	// No-op implementation
}

func (d *DefaultCollector) RecordSuccess(ctx context.Context, metricName string) {
	// No-op implementation
}

// Metric names following required pattern
const (
	// Service metrics
	ServiceVarConfigCreateLatency    = "service_varconfig_create_latency_ms"
	ServiceVarConfigGetLatency       = "service_varconfig_get_latency_ms"
	ServiceVarConfigListLatency      = "service_varconfig_list_latency_ms"
	ServiceVarConfigUpdateLatency     = "service_varconfig_update_latency_ms"
	ServiceVarConfigDeleteLatency     = "service_varconfig_delete_latency_ms"
	
	// Storage metrics
	StorageVarConfigCreateLatency     = "storage_varconfig_create_latency_ms"
	StorageVarConfigGetLatency        = "storage_varconfig_get_latency_ms"
	StorageVarConfigListLatency       = "storage_varconfig_list_latency_ms"
	StorageVarConfigUpdateLatency     = "storage_varconfig_update_latency_ms"
	StorageVarConfigDeleteLatency     = "storage_varconfig_delete_latency_ms"
	StorageVarConfigCreateFailures    = "storage_varconfig_create_failures_total"
	StorageVarConfigGetFailures       = "storage_varconfig_get_failures_total"
	StorageVarConfigListFailures      = "storage_varconfig_list_failures_total"
	StorageVarConfigUpdateFailures     = "storage_varconfig_update_failures_total"
	StorageVarConfigDeleteFailures     = "storage_varconfig_delete_failures_total"
)
