package internal_test

import (
	"bootcamp-web/internal"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWarehouse_ValidName(t *testing.T) {
	name := "Warehouse 1"
	warehouse, err := internal.NewWarehouse(name)
	require.NoError(t, err)
	require.Equal(t, name, warehouse.Name())
}

func TestNewWarehouse_EmptyName(t *testing.T) {
	_, err := internal.NewWarehouse("")
	require.EqualError(t, err, "name must not be empty")
}

func TestNewWarehouse_NoStock(t *testing.T) {
	warehouse, _ := internal.NewWarehouse("Warehouse 1")
	require.Equal(t, 0, warehouse.Stock("Book"))
	require.Equal(t, 0, warehouse.Stock("Pen"))
}

func TestAddStock(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Book", 10)
	require.NoError(t, err)

	err = warehouse.AddStock("Pen", 20)
	require.NoError(t, err)

	require.Equal(t, 10, warehouse.Stock("Book"))
	require.Equal(t, 20, warehouse.Stock("Pen"))
}

func TestAddStock_InvalidQuantity(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Book", -5)
	require.EqualError(t, err, "quantity must be positive")
	require.Equal(t, 0, warehouse.Stock("Book"))
}

func TestAddStock_InvalidProduct(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Invalid", 10)
	require.EqualError(t, err, "invalid product \"Invalid\"")
	require.Equal(t, 0, warehouse.Stock("Invalid"))
}

func TestFulfillOrder(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Book", 10)
	require.NoError(t, err)

	err = warehouse.AddStock("Pen", 20)
	require.NoError(t, err)

	order := internal.Order{
		internal.ProductBook: 5,
		internal.ProductPen:  10,
	}

	err = warehouse.FulfillOrder(order)
	require.NoError(t, err)

	require.Equal(t, 5, warehouse.Stock("Book"))
	require.Equal(t, 10, warehouse.Stock("Pen"))
}

func TestFulfillOrder_NotEnoughStock(t *testing.T) {
	warehouse, err := internal.NewWarehouse("Warehouse 1")
	require.NoError(t, err)

	err = warehouse.AddStock("Book", 10)
	require.NoError(t, err)

	err = warehouse.AddStock("Pen", 20)
	require.NoError(t, err)

	order := internal.Order{
		internal.ProductBook: 15,
		internal.ProductPen:  10,
	}

	err = warehouse.FulfillOrder(order)
	require.EqualError(t, err, "not enough stock for product \"Book\"")
	require.Equal(t, 10, warehouse.Stock("Book"))
	require.Equal(t, 20, warehouse.Stock("Pen"))
}
