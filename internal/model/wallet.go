package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        int64           `gorm:"column:id;primaryKey"` // ID
	UserID    int64           `gorm:"column:user_id"`       // 用戶ID
	Token     decimal.Decimal `gorm:"column:token"`         // 平台幣
	Points    int32           `gorm:"column:points"`        // 平台點數
	CreatedAt time.Time       `gorm:"column:created_at"`    // 創建時間
	UpdatedAt time.Time       `gorm:"column:updated_at"`    // 更新時間
}

type TokenOperation struct {
	Operation NumericOperation
	Token     decimal.Decimal
}

type PointOperation struct {
	Operation NumericOperation
	Points    int32
}
