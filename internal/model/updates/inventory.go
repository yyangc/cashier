package updates

import "cashier/internal/model"

type Inventory struct {
	TotalQuantity     *model.QuantityOperation
	AvailableQuantity *model.QuantityOperation
}
