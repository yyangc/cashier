package updates

import "cashier/internal/model"

type Product struct {
	Quantity          *model.QuantityOperation
	InventoryQuantity *model.QuantityOperation
}
