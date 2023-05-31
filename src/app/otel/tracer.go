package otel

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracerSingleLock = &sync.Mutex{}

var tracerInstance trace.Tracer

func GetTracerInstance() trace.Tracer {
	if tracerInstance == nil {
		tracerSingleLock.Lock()

		defer tracerSingleLock.Unlock()
		if tracerInstance == nil {
      // TODO add service name though config
			tracerInstance = otel.Tracer("")
		}
	}
	return tracerInstance
}
