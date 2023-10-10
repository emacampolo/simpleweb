package internal_test

import (
	"bootcamp-web/internal"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Book", 10)
	require.NoError(t, err)

	warehouseMock := internal.NewWarehouseRepositoryMock()
	warehouseMock.GetFunc = func(ctx context.Context, warehouseID int) (*internal.Warehouse, error) {
		return warehouse, nil
	}

	warehouseMock.UpdateFunc = func(ctx context.Context, warehouse *internal.Warehouse) error {
		return nil
	}

	orderService := internal.NewOrderService(warehouseMock)
	err = orderService.CreateOrder(context.Background(), 1, map[string]int{"Book": 5})
	require.NoError(t, err)

	require.Equal(t, 5, warehouse.Stock("Book"))
}
