package JaegerTracer

import (
	"io"
	"time"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func NewJaegerTracer(serviceName string, jagentHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	jcfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jagentHost,
		},
		ServiceName:serviceName,
	}
	tracer, closer, err = jcfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return
	}

	opentracing.SetGlobalTracer(tracer)
	return
}