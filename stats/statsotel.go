package stats

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var _ Handler = (*OpenTelemetry)(nil)

type (
	OpenTelemetry struct {
		OpenTelemetryOptions
		cacheName string
		counter   metric.Int64Counter
	}

	OpenTelemetryOption func(*OpenTelemetryOptions)

	OpenTelemetryOptions struct {
		namespace string
		subsystem string
		name      string
		help      string
		mp        metric.MeterProvider
	}
)

func OTelWithNamespace(namespace string) OpenTelemetryOption {
	return func(o *OpenTelemetryOptions) {
		o.namespace = namespace
	}
}

func OTelWithSubsystem(subsystem string) OpenTelemetryOption {
	return func(o *OpenTelemetryOptions) {
		o.subsystem = subsystem
	}
}

func OTelWithName(name string) OpenTelemetryOption {
	return func(o *OpenTelemetryOptions) {
		o.name = name
	}
}

func OTelWithHelp(help string) OpenTelemetryOption {
	return func(o *OpenTelemetryOptions) {
		o.help = help
	}
}

func OTelWithMeterProvider(mp metric.MeterProvider) OpenTelemetryOption {
	return func(o *OpenTelemetryOptions) {
		o.mp = mp
	}
}

func NewOpenTelemetry(cacheName string, opts ...OpenTelemetryOption) (*OpenTelemetry, error) {
	o := OpenTelemetryOptions{
		name: "cache_handle_total",
		mp:   otel.GetMeterProvider(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	counter, err := o.mp.
		Meter("jetcache").
		Int64Counter(
			o.name,
			metric.WithDescription(o.help),
		)
	if err != nil {
		return nil, err
	}

	return &OpenTelemetry{
		counter: counter,
	}, nil
}

func (o *OpenTelemetry) IncrHit() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "total"),
			attribute.String("method", "hit"),
		))
}

func (o *OpenTelemetry) IncrMiss() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "total"),
			attribute.String("method", "miss"),
		))
}

func (o *OpenTelemetry) IncrLocalHit() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "local"),
			attribute.String("method", "hit"),
		))
}

func (o *OpenTelemetry) IncrLocalMiss() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "local"),
			attribute.String("method", "miss"),
		))
}

func (o *OpenTelemetry) IncrRemoteHit() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "remote"),
			attribute.String("method", "hit"),
		))
}

func (o *OpenTelemetry) IncrRemoteMiss() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "remote"),
			attribute.String("method", "miss"),
		))
}

func (o *OpenTelemetry) IncrQuery() {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "query"),
			attribute.String("method", "query"),
		))
}

func (o *OpenTelemetry) IncrQueryFail(err error) {
	o.counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("cache_name", o.cacheName),
			attribute.String("cache_type", "query"),
			attribute.String("method", "queryFail"),
			attribute.String("err", err.Error()),
		))
}
