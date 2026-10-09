package telemetry

import (
	"gorm.io/gorm"
)

// struct Telemetry 遥测数据模型
type Telemetry struct {
	gorm.Model

	DeviceID   uint   `gorm:"not null;index;comment:设备ID" json:"device_id"`
	Identifier string `gorm:"type:varchar(128);not null;index;comment:数据类型" json:"identifier"`
	Value      string `gorm:"type:json;not null;comment:数据值" json:"value"`
	Timestamp  int64  `gorm:"not null;index;comment:时间戳" json:"timestamp"`
}

func (Telemetry) TableName() string {
	return "telemetry"
}
