package log

import (
	"context"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
	"wxch-dashboard/config"
)

func InitSentry() {
	debugConfig := config.Get().Debug
	err := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: debugConfig.SentryDSN,

		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		Environment: debugConfig.SentryEnv,

		// Release name
		Release: "wxch-dashboard",

		// Enable printing of SDK debug messages.
		Debug: debugConfig.Verbose,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func FlushSentry(ctx context.Context) {
	if deadline, ok := ctx.Deadline(); ok {
		sentry.Flush(deadline.Sub(time.Now()))
	} else {
		sentry.Flush(3 * time.Second)
	}
}
