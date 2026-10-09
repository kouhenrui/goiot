package app

import (
	"context"
	"errors"
	"fmt"
	"goiot/internal/app/device"
	"goiot/internal/config"
	"goiot/internal/db"
	deviceRep "goiot/internal/domain/device"
	"goiot/internal/handler"
	mqttinf "goiot/internal/infrastructure/mqtt"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"gorm.io/gorm"
)

type App struct {
	db         *gorm.DB
	server     *http.Server
	mqttServer *mqttinf.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	// 初始化 SQLite
	db, err := db.NewSQLite(
		ctx,
		cfg.SQLite.Path,
		cfg.SQLite.Username,
		cfg.SQLite.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("init sqlite: %w", err)
	}
	// 初始化设备仓库和认证服务
	deviceRepo := deviceRep.NewDeviceRepository(db)
	deviceAuth := device.NewAuthService(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceAuth)

	mqttServer, err := mqttinf.NewServer(deviceAuth, cfg.MQTT.Address)
	if err != nil {
		return nil, fmt.Errorf("init mqtt server: %w", err)
	}
	//组装处理器
	handlerGroup := handler.Handlers{
		Device: deviceHandler,
	}
	// 初始化 Gin
	r := handler.New(&handlerGroup)
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		db:         db,
		server:     server,
		mqttServer: mqttServer,
	}, nil

}

func (a *App) Start(ctx context.Context) error {
	glog.Info(ctx, "http server listening on", a.server.Addr)

	go func() {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			glog.Error(ctx, "http server stopped:", err)
		}
		if err = a.mqttServer.Start(); err != nil {
			glog.Error(ctx, "mqtt server stopped:", err)
		}
	}()

	return nil
}
func (a *App) Shutdown(ctx context.Context) error {
	glog.Info(ctx, "shutting down application")

	// 先停止 HTTP Server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf(
			"shutdown http server: %w",
			err,
		)
	}
	glog.Info(ctx, "http server stopped")
	if err := a.mqttServer.Close(); err != nil {
		return fmt.Errorf("shutdown mqtt server: %w", err)
	}
	// 再关闭数据库
	if err := db.Close(a.db); err != nil {
		return fmt.Errorf(
			"close sqlite: %w",
			err,
		)
	}

	glog.Info(ctx, "sqlite closed")

	return nil
}
