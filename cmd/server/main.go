package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"

	"goiot/internal/app"
	"goiot/internal/config"
)

func main() {
	// ============================================================
	// 1. 初始化 Context
	// ============================================================
	ctx := gctx.New()

	cfg, err := config.Load(ctx, "internal/config/config-dev.yaml")
	if err != nil {
		glog.Fatal(ctx, "load config failed:", err)
	}
	if err := config.InitLogger(cfg.Logger); err != nil {
		glog.Fatal(ctx, "init logger failed:", err)
	}
	application, err := app.New(ctx, cfg)
	if err != nil {
		glog.Fatal(ctx, "init app failed:", err)
	}
	if err := application.Start(ctx); err != nil {
		glog.Fatal(ctx, "start app failed:", err)
	}

	// 5. 监听退出信号
	sigCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-sigCtx.Done()

	// 6. 优雅关闭应用，最多等待 10 秒
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		glog.Error(ctx, "shutdown app failed:", err)
		os.Exit(1)
	}

	glog.Info(ctx, "app stopped")
}
