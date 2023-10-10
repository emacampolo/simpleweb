package storagememory

import (
	"bootcamp-web/internal"
	"context"
	"fmt"
)

// warehouseDB represents the data stored in the in-memory repository. The in-memory
// repository does not directly store internal.Warehouse objects, but a
// representation of the warehouse attributes. This makes the in-memory repository
// behave more like a mysql or kvs repository, where you cannot store object
// instances.
type warehouseDB struct {
	id    int
	name  string
	stock map[string]int
}

// WarehouseRepository is an in-memory implementation
// of the internal.WarehouseRepository interface.
type WarehouseRepository struct {
	// lastID is the last id assigned by the repository. It starts
	// at 0, so the first warehouse added to this repository will be 1.
	lastID         int
	warehousesByID map[int]*warehouseDB
}

// NewWarehouseRepository creates a new in-memory warehouse repository.
func NewWarehouseRepository() *WarehouseRepository {
	return &WarehouseRepository{
		warehousesByID: make(map[int]*warehouseDB),
	}
}

// Get implements the internal.WarehouseRepository interface.
func (r *WarehouseRepository) Get(_ context.Context, id int) (*internal.Warehouse, error) {
	warehouseDB, found := r.warehousesByID[id]
	if !found {
		return nil, nil
	}

	stock := make(map[internal.Product]int)
	for product, quantity := range warehouseDB.stock {
		parsedProduct, err := internal.FindProductByName(product)
		if err != nil {
			return nil, err
		}
		stock[parsedProduct] = quantity
	}

	return internal.NewWarehouseFromRepository(warehouseDB.id, internal.WarehouseAttributes{
		Name:  warehouseDB.name,
		Stock: stock,
	}), nil
}

// FindByName implements the internal.WarehouseRepository interface.
func (r *WarehouseRepository) FindByName(ctx context.Context, name string) (*internal.Warehouse, error) {

	for _, warehouseDB := range r.warehousesByID {
		if warehouseDB.name == name {
			return r.Get(ctx, warehouseDB.id)
		}
	}

	return nil, nil
}

// Add implements the internal.WarehouseRepository interface.
func (r *WarehouseRepository) Add(_ context.Context, warehouse *internal.Warehouse) (*internal.Warehouse, error) {
	r.lastID++
	id := r.lastID

	attributes := warehouse.Attributes()

	if warehouse.ID() != 0 {
		return nil, fmt.Errorf("warehouse %q already added with ID %d", warehouse.Name(), warehouse.ID())
	}

	warehouseDB := warehouseDB{
		id:    id,
		name:  warehouse.Name(),
		stock: make(map[string]int),
	}

	for product, quantity := range attributes.Stock {
		warehouseDB.stock[product.Name()] = quantity
	}

	r.warehousesByID[warehouseDB.id] = &warehouseDB
	return internal.NewWarehouseFromRepository(id, attributes), nil
}

// Update implements the internal.WarehouseRepository interface.
func (r *WarehouseRepository) Update(_ context.Context, warehouse *internal.Warehouse) error {
	warehouseDB, found := r.warehousesByID[warehouse.ID()]
	if !found {
		return fmt.Errorf("warehouse %q not found", warehouse.Name())
	}

	warehouseDB.stock = make(map[string]int)
	for product, quantity := range warehouse.Attributes().Stock {
		warehouseDB.stock[product.Name()] = quantity
	}

	return nil
}
