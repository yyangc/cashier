package db

import (
	"context"
	"encoding/json"
	"time"

	"cashier/internal/model"
	"cashier/internal/pkg/errors"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type order struct {
	ID            string          `gorm:"column:id"`
	UserID        int64           `gorm:"column:user_id"`        // 用戶ID
	OriginalPrice decimal.Decimal `gorm:"column:original_price"` // 原始價格
	FinalPrice    decimal.Decimal `gorm:"column:final_price"`    // 最終價格 (扣除優惠活動)
	UsedPoints    int32           `gorm:"column:used_points"`    // 使用平台點數
	PromotionIDs  datatypes.JSON  `gorm:"column:promotion_ids"`  // 使用的優惠ID
	CreatedAt     time.Time       `gorm:"column:created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at"`

	Items []*orderItem `gorm:"foreignKey:OrderID;references:ID"`
}

func (o order) TableName() string {
	return "orders"
}

func (db *database) CreateOrder(ctx context.Context, mOrder *model.Order) (err error) {
	var _order = &order{
		ID:            mOrder.ID,
		UserID:        mOrder.UserID,
		OriginalPrice: mOrder.OriginalPrice,
		FinalPrice:    mOrder.FinalPrice,
		UsedPoints:    mOrder.UsedPoints,
		Items:         make([]*orderItem, 0, len(mOrder.Items)),
	}

	_order.PromotionIDs, err = json.Marshal(mOrder.PromotionIDs)
	if err != nil {
		return errors.Wrapf(errors.ErrInternalServerError, "%+v", err)
	}

	for i := range mOrder.Items {
		_order.Items = append(_order.Items, newOrderItem(mOrder.Items[i]))
	}

	if err := db.WriteDB(ctx).Create(_order).Error; err != nil {
		return errors.Wrapf(duplicateOrInternalError(err), "%+v", err)
	}

	return nil
}

type orderItem struct {
	ID        int64           `gorm:"column:id"`
	OrderID   string          `gorm:"column:order_id"`   // 關聯的 OrderID
	ProductID int64           `gorm:"column:product_id"` // Product 的 ID
	Name      string          `gorm:"column:name"`       // 商品名稱
	UnitPrice decimal.Decimal `gorm:"column:unit_price"` // 平台幣/單價
	Quantity  int32           `gorm:"column:quantity"`   // 數量
}

func (o orderItem) TableName() string {
	return "order_items"
}

func newOrderItem(item *model.OrderItem) *orderItem {
	return &orderItem{
		OrderID:   item.OrderID,
		ProductID: item.ProductID,
		Name:      item.Name,
		UnitPrice: item.UnitPrice,
		Quantity:  item.Quantity,
	}
}
