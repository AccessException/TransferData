package JaegerTracer

import (
	"io"
	"time"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	uberPrometheus "github.com/uber/jaeger-lib/metrics/prometheus"
	clientPrometheus "github.com/prometheus/client_golang/prometheus"
)

func NewJaegerTracer(serviceName string, jagentHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	jcfg := jaegercfg.Configuration{
		//ConstSampler，全量采集
		//ProbabilisticSampler ，概率采集，默认万份之一
		//RateLimitingSampler ，限速采集，每秒只能采集一定量的数据
		//RemotelyControlledSampler ，一种动态采集策略，根据当前系统的访问量调节采集策略

		Sampler: &jaegercfg.SamplerConfig{ 	// 	采样10%
			Type:  "probabilistic",
			Param: 0.1,
		},

		//Sampler: &jaegercfg.SamplerConfig{ 	// 	全量
		//	Type:  "const",
		//	Param: 1,
		//},

		//Sampler: &jaegercfg.SamplerConfig{ 		// 	限速（有待探索）
		//	Type:  "rateLimiting",
		//	Param: 0.5,
		//},

		//withLogSpans -------是否日志上报
		//withMaxQueueSize -------数据最大累计量
		//withFlushInterval -------报告间隔的刷新( ms )
		Reporter: &jaegercfg.ReporterConfig{ // 上报
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			//QueueSize: 2,
			LocalAgentHostPort:  jagentHost,
		},
		ServiceName:serviceName,
		//RPCMetrics:true,
	}

	metricsFactory := uberPrometheus.New(uberPrometheus.WithRegisterer(clientPrometheus.NewPedanticRegistry()))
	tracer, closer, err = jcfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(metricsFactory),
	)
	if err != nil {
		return
	}

	opentracing.SetGlobalTracer(tracer)
	return
}