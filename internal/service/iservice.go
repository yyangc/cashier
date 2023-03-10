package service

import (
	"cashier/internal/model"
	"cashier/internal/model/query"
	"context"
)

type IService interface {
	IOrderService
	IPromotionService
}

type IOrderService interface {
	CreateOrder(ctx context.Context, userID int64, points int32, shoppingCart map[string]int32) (orderID string, err error)
}

type IPromotionService interface {
	ListPromotions(ctx context.Context, promotion query.PromotionOptions) ([]*model.Promotion, error)
}
