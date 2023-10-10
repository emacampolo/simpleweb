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
		var request struct {
			Name string `json:"name" validate:"required"`
		}

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

		type response struct {
			ID int `json:"id"`
		}

		return web.EncodeJSON(w, response{warehouse.ID()}, http.StatusCreated)
	}
}

// NewAddProductStock returns a handler that will add an existing product
// to a warehouse.
func NewAddProductStock(warehouseRepository internal.WarehouseRepository) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var request struct {
			ProductName string `json:"product_name" validate:"required"`
			Quantity    int    `json:"quantity" validate:"required"`
		}

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

		if err := warehouse.AddStock(request.ProductName, request.Quantity); err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		if err := warehouseRepository.Update(r.Context(), warehouse); err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		type response struct {
			ProductName string `json:"product_name"`
			Quantity    int    `json:"quantity"`
		}

		return web.EncodeJSON(w, response{
			ProductName: request.ProductName,
			Quantity:    warehouse.Stock(request.ProductName),
		}, http.StatusOK)
	}
}

// NewCreateOrder returns a handler that will create a new order for a given
// warehouse and a list of products.
func NewCreateOrder(warehouseRepository internal.WarehouseRepository, orderService *internal.OrderService) web.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var request struct {
			Products []struct {
				ProductName string `json:"product_name" validate:"required"`
				Quantity    int    `json:"quantity" validate:"required"`
			} `json:"products" validate:"required"`
		}

		if err := web.DecodeJSON(r, &request); err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid request body: %v", err)
		}

		warehouseID, err := web.ParamInt(r, "warehouse_id")
		if err != nil {
			return NewErrorf(http.StatusBadRequest, "invalid warehouse ID: %v", err)
		}

		products := make(map[string]int)
		for _, product := range request.Products {
			products[product.ProductName] = product.Quantity
		}

		if err := orderService.CreateOrder(r.Context(), warehouseID, products); err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		warehouse, err := warehouseRepository.Get(r.Context(), warehouseID)
		if err != nil {
			return NewError(http.StatusInternalServerError, err.Error())
		}

		var response struct {
			Stock []struct {
				ProductName string `json:"product_name"`
				Quantity    int    `json:"quantity"`
			} `json:"stock"`
		}

		// For every product in the order, get the current stock
		// and add it to the response.
		for productName := range products {
			response.Stock = append(response.Stock, struct {
				ProductName string `json:"product_name"`
				Quantity    int    `json:"quantity"`
			}{
				ProductName: productName,
				Quantity:    warehouse.Stock(productName),
			})
		}

		return web.EncodeJSON(w, response, http.StatusCreated)
	}
}
