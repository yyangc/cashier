package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductStatus int8

const (
	ProductStatusUnknown ProductStatus = iota
	ProductStatusOn                    // 上架
	ProductStatusDown                  // 下架
)

func (p ProductStatus) Str() string {
	switch p {
	case ProductStatusOn:
		return "On"
	case ProductStatusDown:
		return "Down"
	default:
		return "Unknown"
	}
}

// Product 產品
type Product struct {
	ID        int64
	Name      string          // 商品名稱
	Status    ProductStatus   // 商品上下架狀態
	Price     decimal.Decimal // 價格(單位：平台幣)
	CreatedAt time.Time       // 創建時間
	UpdatedAt time.Time       // 更新時間

	Inventory *Inventory // 庫存
}

func (p *Product) IsAvailable() bool {
	if p.Status == ProductStatusOn {
		return true
	}
	return false
}
