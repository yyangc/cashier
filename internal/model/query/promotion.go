package query

import (
	"cashier/internal/model"

	"time"
)

type PromotionOptions struct {
	IDIn       []int64               // 活動ID
	TypeIn     []model.PromotionType // 活動類型
	StartAtGte *time.Time            // 活動時間大於等於
	EndAtLt    *time.Time            // 活動時間小於
}
