package device

import "gorm.io/gorm"

type ConnectionType string

const (
	ConnectionSerial ConnectionType = "serial"
	ConnectionMQTT   ConnectionType = "mqtt"
)

type Status string

const (
	StatusInactive Status = "inactive"
	StatusOnline   Status = "online"
	StatusOffline  Status = "offline"
	StatusDisabled Status = "disabled"
)

// struct Device 设备模型
type Device struct {
	gorm.Model

	ProductID        uint           `gorm:"not null;index;comment:产品ID" json:"product_id"`
	Name             string         `gorm:"type:varchar(128);not null;comment:设备名称" json:"name"`
	DeviceKey        string         `gorm:"type:varchar(128);not null;uniqueIndex;comment:设备密钥" json:"device_key"`
	ConnectionType   ConnectionType `gorm:"type:varchar(16);not null;index;comment:连接类型" json:"connection_type"`
	ConnectionConfig string         `gorm:"type:text;comment:连接配置" json:"connection_config"`

	// 仅存 bcrypt 哈希，禁止通过 JSON 返回。
	SecretHash string `gorm:"type:varchar(128);comment:设备密钥哈希" json:"-"`
	Status     Status `gorm:"type:varchar(32);not null;default:inactive;index;comment:设备状态" json:"status"`
}

func (Device) TableName() string {
	return "devices"
}
