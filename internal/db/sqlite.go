package db

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"goiot/internal/domain/device"
	"goiot/internal/domain/product"
	"goiot/internal/domain/telemetry"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewSQLite(ctx context.Context, path, username, password string) (*gorm.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite path is empty")
	}
	if username == "" || password == "" {
		return nil, fmt.Errorf("sqlite username and password must be configured")
	}

	// 创建数据库所在目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf(
			"create sqlite directory %q: %w",
			dir,
			err,
		)
	}
	dsn := "file:" + filepath.ToSlash(path) + "?" + url.Values{
		"_auth":       {""},
		"_auth_user":  {username},
		"_auth_pass":  {password},
		"_auth_crypt": {"sha256"},
	}.Encode()
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	if err := automigrate(db); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("automigrate sqlite: %w", err)
	}

	return db, nil
}

func automigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&device.Device{}, &product.Product{}, &telemetry.Telemetry{}); err != nil {
		return fmt.Errorf("auto migrate device: %w", err)
	}
	return nil
}

// Close 关闭 SQLite 数据库。

func Close(db *gorm.DB) error {

	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf(
			"get sqlite sql db: %w",
			err,
		)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf(
			"close sqlite: %w",
			err,
		)
	}
	return nil

}
