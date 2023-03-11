package model

import "time"

// Inventory 庫存
type Inventory struct {
	ID                int64
	ProductID         int64 // 關聯 Product.ID
	TotalQuantity     int32 // 總庫存數量
	AvailableQuantity int32 // 可售庫存數量
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type QuantityOperation struct {
	Operation NumericOperation
	Quantity  int32
}
