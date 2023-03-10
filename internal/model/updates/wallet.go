package updates

import (
	"cashier/internal/model"
)

type Wallet struct {
	TokenOperation  *model.TokenOperation // 平台幣操作
	PointsOperation *model.PointOperation // 平台點數操作

	Lock bool // lock for update
}
