package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID            string
	UserID        int64           // 用戶ID
	OriginalPrice decimal.Decimal // 原始價格
	FinalPrice    decimal.Decimal // 最終價格 (扣除優惠活動)
	UsedPoints    int32           // 使用平台點數
	PromotionIDs  []int64         // 使用的優惠ID
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Items      []*OrderItem // 關聯的商品 Product
	Promotions []*Promotion // 使用的優惠 Promotion
}

// OrderItem 訂單詳情的紀錄
type OrderItem struct {
	ID        int64
	OrderID   string          // 關聯的 OrderID
	ProductID string          // Product 的 ID
	Name      string          // 商品名稱
	UnitPrice decimal.Decimal // 平台幣/單價
	Quantity  int32           // 數量
}

type PurchaseProduct struct {
	Product  *Product
	Quantity int32 // 數量
}

func NewOrderItem(orderID string, product *Product, quantity int32) *OrderItem {
	return &OrderItem{
		OrderID:   orderID,
		ProductID: product.ID,
		Name:      product.Name,
		UnitPrice: product.Price,
		Quantity:  quantity,
	}
}
