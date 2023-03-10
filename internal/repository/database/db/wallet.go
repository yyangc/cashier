package db

import (
	"cashier/internal/model"
	"cashier/internal/model/query"
	"cashier/internal/model/updates"
	"context"
)

func (db *database) GetWallet(ctx context.Context, options *query.WalletOptions) (*model.Wallet, error) {
	// skip
	return &model.Wallet{}, nil
}

func (db *database) UpdateWallet(ctx context.Context, options *query.WalletOptions, updates *updates.Wallet) error {
	// skip
	return nil
}
