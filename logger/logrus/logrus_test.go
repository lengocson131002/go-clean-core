package logrus

import (
	"context"
	"testing"
	"time"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/go-clean-core/trace"
	"github.com/lengocson131002/go-clean-core/trace/otel"
)

func TestMaskedData(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          `{"username": "john_doe", "password": "123"} <password> 1234 </password> <credentials> 12345 </credentials> base64data: ZnNkZnNkZnNkZnNkZnNkZnNkZnNkZnNmc2RzZGZzZGZkc2ZzZGZzZGZzZGZmc2QK`,
			expectedOutput: `{"username": "john_doe", "password": "***"} <password> **** </password> <credentials> ***** </credentials> base64data: ****************************************************************`,
		},
		{
			input:          `{"password": "123"} <password>1234</password> <credentials>12345</credentials> base64data: XYZ123==`,
			expectedOutput: `{"password": "***"} <password>****</password> <credentials>*****</credentials> base64data: XYZ123==`,
		},
		{
			input:          `<![CDATA[ <password> 123 </password> ]]>, <![CDATA[ <credentials> user:1234 </credentials> ]]>, base64data: MNO456==`,
			expectedOutput: `<![CDATA[ <password> *** </password> ]]>, <![CDATA[ <credentials> ********* </credentials> ]]>, base64data: MNO456==`,
		},
		{
			input:          `base64data: ABCDEFGH12345==, <password>123</password>, {"password": "1234"}, <credentials>user:1234</credentials>`,
			expectedOutput: `base64data: ABCDEFGH12345==, <password>***</password>, {"password": "****"}, <credentials>*********</credentials>`,
		},
		{
			input:          `nothing`,
			expectedOutput: `nothing`,
		},
		// Add more test cases as needed
	}

	logger := NewLogrusLogger()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			logger.Info(ctx, tc.input)
		})
	}
}

func TestLogCommon(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "log information 1",
		}, {
			input: "log information 2",
		},
	}

	ctx := context.Background()
	logger := NewLogrusLogger()

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			logger.Info(ctx, tc.input)
		})
	}
}

func TestLogTracing(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "log information 1",
		}, {
			input: "log information 2",
		},
	}

	tracer, err := otel.NewOpenTelemetryTracer(context.Background(),
		trace.WithTraceServiceName("go_logrus_testing"),
		trace.WithTraceExporterEndpoint("localhost:4318"))

	if err != nil {
		t.Error(err)
	}

	logger := NewLogrusLogger(
		logger.WithTracer(tracer),
	)

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			ctx := context.Background()
			ctx, finish := tracer.StartInternalTrace(ctx, "testing")
			logger.Info(ctx, tc.input)
			finish(ctx, trace.WithInternalResponse("test"))
		})
	}

	time.Sleep(5 * time.Second)
}
