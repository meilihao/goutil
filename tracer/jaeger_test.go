package tracer

import (
	"testing"

	"go.uber.org/zap"
)

func TestTimeSince(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	conf := &JaegerConfig{
		Name:      "test",
		Enable:    true,
		AgentAddr: "127.0.0.1:6831",
		LogSpans:  true,
	}

	tracer, err := InitTracer(conf, sugar)
	if err != nil {
		t.Fatal(err)
	}

	tracer.Close()
}
