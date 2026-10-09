package product

import "gorm.io/gorm"

// struct Product 产品模型
type Product struct {
	gorm.Model

	Name        string `gorm:"type:varchar(128);not null;index;comment:产品名称" json:"name"`
	Description string `gorm:"type:text;comment:产品描述" json:"description"`
	Protocol    string `gorm:"type:varchar(32);not null;comment:产品协议" json:"protocol"`
	Status      string `gorm:"type:varchar(32);not null;default:enabled;index;comment:产品状态" json:"status"`
}
