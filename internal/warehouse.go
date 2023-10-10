package internal

import (
	"context"
	"fmt"
)

// Warehouse represents a warehouse of products.
type Warehouse struct {
	id         int
	attributes WarehouseAttributes
}

// WarehouseAttributes represents the attributes of a warehouse
// that can be stored in a repository.
type WarehouseAttributes struct {
	Name  string
	Stock map[Product]int
}

// NewWarehouse creates a new warehouse. The name must not be empty.
func NewWarehouse(name string) (*Warehouse, error) {
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	return &Warehouse{
		attributes: WarehouseAttributes{
			Name:  name,
			Stock: make(map[Product]int),
		},
	}, nil
}

// NewWarehouseFromRepository creates a new warehouse.
// It must be used by repositories to create a warehouse.
func NewWarehouseFromRepository(id int, attributes WarehouseAttributes) *Warehouse {
	return &Warehouse{
		id:         id,
		attributes: attributes,
	}
}

// Name returns the name of the warehouse, never empty.
func (w *Warehouse) Name() string {
	return w.attributes.Name
}

// ID returns the ID of the warehouse.
// It returns 0 if the warehouse has not been stored in a repository.
func (w *Warehouse) ID() int {
	return w.id
}

// AddStock adds stock of a product to the warehouse.
// It returns an error if the product is invalid or if the quantity is negative.
func (w *Warehouse) AddStock(productName string, quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("quantity must be positive")
	}

	product, err := FindProductByName(productName)
	if err != nil {
		return err
	}

	w.attributes.Stock[product] += quantity
	return nil
}

// FulfillOrder fulfills an order from the warehouse.
// It returns an error if the warehouse does not have enough stock
// to fulfill the order.
// The operation is atomic meaning that if the warehouse does not have
// enough stock to fulfill the order, the warehouse stock is not changed.
func (w *Warehouse) FulfillOrder(order Order) error {
	for product, quantity := range order {
		if w.attributes.Stock[product] < quantity {
			return fmt.Errorf("not enough stock for product %q", product.name)
		}
	}

	for product, quantity := range order {
		w.attributes.Stock[product] -= quantity
	}

	return nil
}

// Stock returns the stock of the warehouse
// for a particular product.
func (w *Warehouse) Stock(productName string) int {
	product, err := FindProductByName(productName)
	if err != nil {
		return 0
	}
	return w.attributes.Stock[product]
}

// Attributes returns the attributes of the warehouse.
func (w *Warehouse) Attributes() WarehouseAttributes {
	return w.attributes
}

// WarehouseRepository represents a repository of warehouses.
type WarehouseRepository interface {
	// FindByName finds a warehouse by its name.
	// It returns nil if the warehouse does not exist.
	FindByName(ctx context.Context, name string) (*Warehouse, error)
	// Add adds a warehouse to the repository.
	// It returns an error if the warehouse already exists.
	Add(ctx context.Context, warehouse *Warehouse) (*Warehouse, error)
	// Update updates a warehouse in the repository.
	// It returns an error if the warehouse does not exist.
	Update(ctx context.Context, warehouse *Warehouse) error
	// Get gets a warehouse from the repository by its ID.
	// It returns nil if the warehouse does not exist.
	Get(ctx context.Context, id int) (*Warehouse, error)
}
