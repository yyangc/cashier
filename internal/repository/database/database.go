package database

import (
	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	"context"
)

type IDatabase interface {
	Begin(ctx context.Context) IDatabase
	Commit() error
	Rollback() error
	Transaction(ctx context.Context, callback func(ctx context.Context, txRepo IDatabase) error) error

	IProductDB
	IPromotionDB
	IMemberDB
	IOrderDB
	IWalletDB
	IInventoryDB
}

type IPromotionDB interface {
	// ListPromotions 取得多筆優惠活動
	ListPromotions(ctx context.Context, options *query.PromotionOptions) ([]*model.Promotion, error)
	CreatePromotion(ctx context.Context, mPromotion *model.Promotion) error
}

type IProductDB interface {
	ListProducts(ctx context.Context, options *query.ProductOptions) ([]*model.Product, error)
}

type IMemberDB interface {
	// GetMember 取得用戶的會員等級
	GetMember(ctx context.Context, options *query.MemberOptions) (*model.Member, error)
}

type IOrderDB interface {
	// CreateOrder 建立訂單
	CreateOrder(ctx context.Context, order *model.Order) error
}

type IWalletDB interface {
	// GetWallet 取得用戶的錢包
	GetWallet(ctx context.Context, options *query.WalletOptions) (*model.Wallet, error)
	// UpdateWallet 更新用戶的錢包
	UpdateWallet(ctx context.Context, options *query.WalletOptions, updates *updates.Wallet) error
}

type IInventoryDB interface {
	ListInventories(ctx context.Context, options *query.InventoryOptions) ([]*model.Inventory, error)
	UpdateInventory(ctx context.Context, options *query.InventoryOptions, updates *updates.Inventory) error
}
