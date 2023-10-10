package internal

import (
	"context"
	"fmt"
)

// WarehouseRepositoryMock is a mock implementation of internal.WarehouseRepository.
type WarehouseRepositoryMock struct {
	FindByNameFunc func(ctx context.Context, name string) (*Warehouse, error)
	AddFunc        func(ctx context.Context, warehouse *Warehouse) (*Warehouse, error)
	UpdateFunc     func(ctx context.Context, warehouse *Warehouse) error
	GetFunc        func(ctx context.Context, warehouseID int) (*Warehouse, error)
}

func (w *WarehouseRepositoryMock) FindByName(ctx context.Context, name string) (*Warehouse, error) {
	return w.FindByNameFunc(ctx, name)
}

func (w *WarehouseRepositoryMock) Add(ctx context.Context, warehouse *Warehouse) (*Warehouse, error) {
	return w.AddFunc(ctx, warehouse)
}

func (w *WarehouseRepositoryMock) Update(ctx context.Context, warehouse *Warehouse) error {
	return w.UpdateFunc(ctx, warehouse)
}

func (w *WarehouseRepositoryMock) Get(ctx context.Context, warehouseID int) (*Warehouse, error) {
	return w.GetFunc(ctx, warehouseID)
}

// NewWarehouseRepositoryMock creates a new WarehouseRepositoryMock.
func NewWarehouseRepositoryMock() *WarehouseRepositoryMock {
	return &WarehouseRepositoryMock{
		FindByNameFunc: func(ctx context.Context, name string) (*Warehouse, error) {
			return nil, fmt.Errorf("not implemented")
		},
		AddFunc: func(ctx context.Context, warehouse *Warehouse) (*Warehouse, error) {
			return nil, fmt.Errorf("not implemented")
		},
		UpdateFunc: func(ctx context.Context, warehouse *Warehouse) error {
			return fmt.Errorf("not implemented")
		},
		GetFunc: func(ctx context.Context, warehouseID int) (*Warehouse, error) {
			return nil, fmt.Errorf("not implemented")
		},
	}
}
