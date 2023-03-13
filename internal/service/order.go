package service

import (
	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	"cashier/internal/pkg/errors"
	iDB "cashier/internal/repository/database"
	"context"

	"github.com/rs/xid"

	"github.com/shopspring/decimal"
)

// CreateOrder 建立訂單，返回最終訂單金額
func (s *service) CreateOrder(ctx context.Context, userID int64, points int32, shoppingCart map[int64]int32) (orderID string, err error) {
	order := &model.Order{
		ID:         xid.New().String(),
		UserID:     userID,
		UsedPoints: points,
	}

	// 清算購物車 (取得原始總金額 & 商品清單)
	var products []*model.Product
	order.OriginalPrice, products, err = s.CalculateShoppingCart(ctx, shoppingCart)
	if err != nil {
		return "", err
	}

	// 紀錄該訂單關聯的商品
	var productIDs = make([]int64, 0, len(products))
	for _, product := range products {
		order.Items = append(order.Items, model.NewOrderItem(order.ID, product, shoppingCart[product.ID]))
		productIDs = append(productIDs, product.ID)
	}

	// 返回符合條件的優惠 & 優惠後的訂單金額
	order.FinalPrice, order.PromotionIDs, err = s.CalculateDiscountPrice(ctx, order)
	if err != nil {
		return "", err
	}

	//  建立訂單
	err = s.db.Transaction(ctx, func(txCtx context.Context, txRepo iDB.IDatabase) error {
		// 取得用戶錢包
		wallet, err := txRepo.GetWallet(txCtx, &query.WalletOptions{
			UserIDIn: []int64{userID},
			Lock:     true,
		})
		if err != nil {
			return err
		}

		var updatesWallet = updates.Wallet{}

		// 檢查錢包的平台幣餘額是否足夠
		if order.FinalPrice.GreaterThan(wallet.Token) {
			return errors.Wrapf(errors.ErrInsufficientBalance,
				"Insufficient token, order token %d is greater than wallet token %d", order.FinalPrice, wallet.Token,
			)
		}

		updatesWallet.TokenOperation = &model.TokenOperation{
			Operation: model.NumericOperationSub,
			Token:     order.FinalPrice,
		}

		// 訂單有使用到平台點數，則檢查
		if order.UsedPoints > 0 {
			if order.UsedPoints > wallet.Points {
				return errors.Wrapf(errors.ErrInsufficientBalance,
					"Insufficient point, order point %d is greater than wallet point %d", order.UsedPoints, wallet.Points,
				)
			}
			updatesWallet.PointsOperation = &model.PointOperation{
				Operation: model.NumericOperationSub,
				Points:    order.UsedPoints,
			}
		}

		// 更新用戶錢包 (扣錢 + 扣點數)
		if err := txRepo.UpdateWallet(txCtx,
			&query.WalletOptions{IDIn: []int64{wallet.ID}},
			&updatesWallet,
		); err != nil {
			return err
		}

		// 檢查商品庫存
		inventories, err := txRepo.ListInventories(txCtx, &query.InventoryOptions{
			ProductIDIn: productIDs,
			Lock:        true,
		})
		if err != nil {
			return err
		}

		for i := range inventories {
			// 商品庫存必須大於等於購買數量
			if inventories[i].AvailableQuantity < shoppingCart[products[i].ID] {
				return errors.Wrapf(errors.ErrInsufficientBalance,
					"productID(%d) is out of stock. %d < %d", inventories[i].AvailableQuantity, shoppingCart[products[i].ID],
				)
			}
		}

		// 更新庫存
		for i := range order.Items {
			if err := txRepo.UpdateInventory(txCtx,
				&query.InventoryOptions{ProductIDIn: []int64{order.Items[i].ProductID}},
				&updates.Inventory{AvailableQuantity: &model.QuantityOperation{
					Operation: model.NumericOperationSub,
					Quantity:  order.Items[i].Quantity,
				}},
			); err != nil {
				return err
			}
		}

		// 建立訂單 & 訂單詳情
		if err := txRepo.CreateOrder(txCtx, order); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return order.ID, nil
}

// CalculateShoppingCart 清算購物車的商品，返回總金額 & 商品
func (s *service) CalculateShoppingCart(ctx context.Context, purchaseList map[int64]int32) (
	price decimal.Decimal, products []*model.Product, err error,
) {

	var productIDs = make([]int64, 0, len(purchaseList))
	for id := range purchaseList {
		productIDs = append(productIDs, id)
	}

	products, err = s.db.ListProducts(ctx, &query.ProductOptions{IDIn: productIDs})
	if err != nil {
		return decimal.Zero, nil, err
	}

	if len(products) == 0 {
		return decimal.Zero, nil, errors.WithStack(errors.ErrResourceNotFound)
	}

	var originalPrice decimal.Decimal

	for _, product := range products {
		if product.Status != model.ProductStatusOn {
			return decimal.Zero, nil, errors.Wrapf(errors.ErrResourceUnavailable, "product(%d) status is %s", product.ID, product.Status.Str())
		}
		if product.Inventory.AvailableQuantity <= 0 {
			return decimal.Zero, nil, errors.Wrapf(errors.ErrResourceUnavailable, "product(%d) is sold out", product.ID)
		}

		quantity := purchaseList[product.ID]
		originalPrice = originalPrice.Add(product.Price.Mul(decimal.NewFromInt32(quantity)))
	}

	return originalPrice, products, nil
}

// CalculateDiscountPrice 返回訂單折扣後金額 & 該訂單使用的優惠ID
func (s *service) CalculateDiscountPrice(ctx context.Context, order *model.Order) (afterPrice decimal.Decimal, promotionIDs []int64, err error) {
	// 取得用戶的會員等級
	member, err := s.db.GetMember(ctx, &query.MemberOptions{IDIn: []int64{order.UserID}})
	if err != nil {
		return decimal.Zero, nil, err
	}

	// 取得當前的優惠活動
	promotionMap, err := s.GetCurrPromotionsMap(ctx)
	if err != nil {
		return decimal.Zero, nil, err
	}

	// 依優惠活動計算訂單金額 & 紀錄使用的優惠
	promotionIDs = make([]int64, 0)
	afterPrice = order.OriginalPrice
	calPriceInput := &model.CalculatePriceInput{
		Member:     member,
		UsedPoints: order.UsedPoints,
	}

	for _, pType := range model.ValidPromotionTypes {
		if promotion, exist := promotionMap[pType]; exist {
			var usedPromotion bool
			usedPromotion, afterPrice = promotion.Extension.CalculatePrice(afterPrice, calPriceInput)
			if usedPromotion {
				// 紀錄這個訂單有用到的優惠
				promotionIDs = append(promotionIDs, promotion.ID)
			}
		}
	}

	return afterPrice, promotionIDs, nil
}
