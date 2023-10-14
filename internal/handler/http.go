package handler

import (
	"bootcamp-web/internal"
	"bootcamp-web/platform/web"
	"net/http"
)

// NewAddWarehouse returns a handler that will create a new warehouse
// with a given name.
func NewAddWarehouse(warehouseRepository internal.WarehouseRepository) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var request PostWarehousesJSONRequestBody

		if err := web.DecodeJSON(r, &request); err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid request body: %v", err)
		}

		warehouse, err := internal.NewWarehouse(request.Name)
		if err != nil {
			return NewErrorf(http.StatusBadRequest, err.Error())
		}

		warehouse, err = warehouseRepository.Add(r.Context(), warehouse)
		if err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		return web.EncodeJSON(w, Warehouse{warehouse.ID()}, http.StatusCreated)
	}
}

// NewAddProductToWarehouse returns a handler that will add an existing product
// to a warehouse.
func NewAddProductToWarehouse(warehouseRepository internal.WarehouseRepository) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var request PostWarehousesWarehouseIdProductsJSONRequestBody
		if err := web.DecodeJSON(r, &request); err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid request body: %v", err)
		}

		warehouseID, err := web.ParamInt(r, "warehouse_id")
		if err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid warehouse ID: %v", err)
		}

		warehouse, err := warehouseRepository.Get(r.Context(), warehouseID)
		if err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		if warehouse == nil {
			return NewErrorf(http.StatusNotFound, "warehouse %d not found", warehouseID)
		}

		if err := warehouse.AddStock(request.Product.Name, request.Quantity); err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		if err := warehouseRepository.Update(r.Context(), warehouse); err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		return web.EncodeJSON(w, ProductQuantity{
			Product:  Product{Name: request.Product.Name},
			Quantity: warehouse.Stock(request.Product.Name),
		}, http.StatusOK)
	}
}

// NewCreateOrder returns a handler that will create a new order for a given
// warehouse and a list of products.
func NewCreateOrder(warehouseRepository internal.WarehouseRepository, orderService *internal.OrderService) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var request PostWarehousesWarehouseIdOrdersJSONRequestBody

		if err := web.DecodeJSON(r, &request); err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid request body: %v", err)
		}

		warehouseID, err := web.ParamInt(r, "warehouse_id")
		if err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid warehouse ID: %v", err)
		}

		products := make(map[string]int)
		for _, product := range request.Order {
			products[product.Product.Name] = product.Quantity
		}

		if err := orderService.CreateOrder(r.Context(), warehouseID, products); err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		warehouse, err := warehouseRepository.Get(r.Context(), warehouseID)
		if err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		var response Stock

		// For every product in the order, get the current stock
		// and add it to the response.
		for productName := range products {
			response.Stock = append(response.Stock, ProductQuantity{
				Product:  Product{Name: productName},
				Quantity: warehouse.Stock(productName),
			})
		}

		return web.EncodeJSON(w, response, http.StatusCreated)
	}
}
