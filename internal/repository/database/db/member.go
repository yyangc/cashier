package db

import (
	"context"

	"cashier/internal/model"
	"cashier/internal/model/query"
)

// GetMember 取得該用戶的會員方案
func (db *database) GetMember(ctx context.Context, options *query.MemberOptions) (*model.Member, error) {
	// skip ...
	return &model.Member{}, nil
}
