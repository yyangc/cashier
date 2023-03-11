package db

import (
	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	"context"
)

// GetWallet 取得用戶錢包
func (db *database) GetWallet(ctx context.Context, options *query.WalletOptions) (*model.Wallet, error) {
	// skip
	return &model.Wallet{}, nil
}

// UpdateWallet 更新用戶錢包
func (db *database) UpdateWallet(ctx context.Context, options *query.WalletOptions, updates *updates.Wallet) error {
	// skip
	return nil
}
