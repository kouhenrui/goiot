package device

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(
	ctx context.Context,
	d *Device,
) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *DeviceRepository) FindByKey(
	ctx context.Context,
	key string,
) (*Device, error) {
	var d Device

	err := r.db.WithContext(ctx).
		Where("device_key = ?", key).
		First(&d).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("device not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	return &d, nil
}
