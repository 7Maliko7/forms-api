package oc

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"

	"go.opencensus.io/trace"
)

type BalancerType string

const (
	Random     BalancerType = "random"
	RoundRobin BalancerType = "round robin"
)

func ClientEndpoint(operationName string, attrs ...trace.Attribute) endpoint.Middleware {
	attrs = append(
		attrs, trace.StringAttribute("gokit.endpoint.type", "client"),
	)
	return kitoc.TraceEndpoint(
		"gokit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(attrs...),
	)
}

func ServerEndpoint(operationName string, attrs ...trace.Attribute) endpoint.Middleware {
	attrs = append(
		attrs, trace.StringAttribute("gokit.endpoint.type", "server"),
	)
	return kitoc.TraceEndpoint(
		"gokit/endpoint "+operationName,
		kitoc.WithEndpointAttributes(attrs...),
	)
}

func RetryEndpoint(
	operationName string, balancer BalancerType, max int, timeout time.Duration,
) endpoint.Middleware {
	return kitoc.TraceEndpoint("gokit/retry "+operationName,
		kitoc.WithEndpointAttributes(
			trace.StringAttribute("gokit.balancer.type", string(balancer)),
			trace.StringAttribute("gokit.retry.timeout", timeout.String()),
			trace.Int64Attribute("gokit.retry.max_count", int64(max)),
		),
	)
}
