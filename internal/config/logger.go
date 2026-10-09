package config

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/glog"
)

func InitLogger(cfg Logger) error {
	// 设置日志目录
	fmt.Print(cfg, "9999")
	if err := glog.SetPath(cfg.Path); err != nil {
		return fmt.Errorf("set logger path: %w", err)
	}

	// 设置日志文件名模板
	glog.SetFile(cfg.File)

	// 设置日志级别
	if err := glog.SetLevelStr(cfg.Level); err != nil {
		return fmt.Errorf("set logger level: %w", err)
	}

	// 设置时间格式和终端输出
	glog.SetConfigWithMap(map[string]any{
		"timeFormat":         cfg.TimeFormat,
		"stdout":             cfg.Stdout,
		"rotateSize":         int64(100 * 1024 * 1024),
		"rotateBackupExpire": 7 * 24 * time.Hour,
	})

	return nil
}
