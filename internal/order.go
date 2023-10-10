package internal

import (
	"context"
	"fmt"
)

// Order represents an order of products that
// it fulfilled by a warehouse.
// The key is the product and the value is the quantity.
type Order map[Product]int

// OrderService is the domain service for orders.
type OrderService struct {
	warehouseRepository WarehouseRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(warehouseRepository WarehouseRepository) *OrderService {
	return &OrderService{
		warehouseRepository: warehouseRepository,
	}
}

// CreateOrder creates a new order.
func (o *OrderService) CreateOrder(ctx context.Context, warehouseID int, products map[string]int) error {
	order := make(map[Product]int, len(products))
	for name, quantity := range products {
		product, err := FindProductByName(name)
		if err != nil {
			return err
		}

		order[product] = quantity
	}

	warehouse, err := o.warehouseRepository.Get(ctx, warehouseID)
	if err != nil {
		return err
	}

	if warehouse == nil {
		return fmt.Errorf("warehouse %d not found", warehouseID)
	}

	if err := warehouse.FulfillOrder(order); err != nil {
		return err
	}

	if err := o.warehouseRepository.Update(ctx, warehouse); err != nil {
		return err
	}

	return nil
}
