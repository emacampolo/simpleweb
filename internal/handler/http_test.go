package handler_test

import (
	"bootcamp-web/internal"
	"bootcamp-web/internal/handler"
	"bootcamp-web/platform/web"
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAddWarehouse(t *testing.T) {
	warehouseRepositoryMock := internal.NewWarehouseRepositoryMock()
	warehouseRepositoryMock.AddFunc = func(
		ctx context.Context, warehouse *internal.Warehouse) (*internal.Warehouse, error) {
		return internal.NewWarehouseFromRepository(1, warehouse.Attributes()), nil
	}

	req := httptest.NewRequest("", "/", strings.NewReader(`{"name":"Warehouse 1"}`))
	w := httptest.NewRecorder()

	h := handler.NewAddWarehouse(warehouseRepositoryMock)
	err := h(w, req)
	require.NoError(t, err)

	require.Equal(t, 201, w.Code)
	require.Equal(t, `{"id":1}`, w.Body.String())
}

func TestNewAddProductStock(t *testing.T) {
	warehouseRepositoryMock := internal.NewWarehouseRepositoryMock()
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	warehouseRepositoryMock.UpdateFunc = func(ctx context.Context, warehouse2 *internal.Warehouse) error {
		if warehouse != warehouse2 {
			t.Error("Warehouse should be the same")
		}
		return nil
	}

	warehouseRepositoryMock.GetFunc = func(ctx context.Context, warehouseID int) (*internal.Warehouse, error) {
		if warehouseID == 1 {
			return warehouse, nil
		}

		return nil, nil
	}

	req := httptest.NewRequest("", "/", strings.NewReader(`{"product_name":"Book","quantity":10}`))
	req = web.WithURLParams(t, req, map[string]string{"warehouse_id": "1"})
	w := httptest.NewRecorder()

	h := handler.NewAddProductStock(warehouseRepositoryMock)
	err = h(w, req)
	require.NoError(t, err)

	require.Equal(t, 200, w.Code)
	require.JSONEq(t, `{"product_name":"Book","quantity":10}`, w.Body.String())
	require.Equal(t, 10, warehouse.Stock("Book"))
}

func TestNewCreateOrder(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Book", 10)
	require.NoError(t, err)

	warehouseRepositoryMock := internal.NewWarehouseRepositoryMock()
	warehouseRepositoryMock.GetFunc = func(ctx context.Context, warehouseID int) (*internal.Warehouse, error) {
		return warehouse, nil
	}

	warehouseRepositoryMock.UpdateFunc = func(ctx context.Context, warehouse2 *internal.Warehouse) error {
		if warehouse != warehouse2 {
			t.Error("Warehouse should be the same")
		}
		return nil
	}

	orderService := internal.NewOrderService(warehouseRepositoryMock)

	req := httptest.NewRequest("", "/",
		strings.NewReader(`{"products":[{"product_name":"Book","quantity":5}]}`))
	req = web.WithURLParams(t, req, map[string]string{"warehouse_id": "1"})
	w := httptest.NewRecorder()

	h := handler.NewCreateOrder(warehouseRepositoryMock, orderService)
	err = h(w, req)
	require.NoError(t, err)

	require.Equal(t, 201, w.Code)
	require.JSONEq(t, `{"stock":[{"product_name":"Book","quantity":5}]}`, w.Body.String())
	require.Equal(t, 5, warehouse.Stock("Book"))
}
