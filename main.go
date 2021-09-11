package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wxch-dashboard/config"
	"wxch-dashboard/logic/daemon"
	"wxch-dashboard/logic/log"
	"wxch-dashboard/logic/rpc"
)

func main() {
	go start()

	// gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.GetLogger().Info("gracefully shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	wg := sync.WaitGroup{}
	groupWait := func(cbs ...func()) {
		for _, cb := range cbs {
			curFunc := cb
			wg.Add(1)
			go func() {
				curFunc()
				wg.Done()
			}()
		}

		wg.Wait()
	}

	groupWait(
		func() {
			rpc.Stop(ctx)
			log.GetLogger().Info("rpc server stoped")
		},
		func() {
			daemon.StopCron(ctx)
			log.GetLogger().Info("cron stoped")
		},
		func() {
			_ = log.GetLogger().Sync()
		},
		func() {
			log.FlushSentry(ctx)
			log.GetLogger().Info("sentry flushed")
		},
	)

	log.GetLogger().Info("gracefully shutdown success")
}

func start() {
	// init cron task
	if !config.Get().Debug.DisableCron {
		daemon.InitCron()
	}

	// init sentry
	if !config.Get().Debug.DisableSentry {
		log.InitSentry()
	}

	// launch RPC server
	if err := rpc.StartAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
