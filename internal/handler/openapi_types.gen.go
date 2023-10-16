// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package handler

// AddProductToWarehouse defines model for AddProductToWarehouse.
type AddProductToWarehouse struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

// AddWarehouse defines model for AddWarehouse.
type AddWarehouse struct {
	Name string `json:"name"`
}

// CreateOrder defines model for CreateOrder.
type CreateOrder struct {
	Order []ProductQuantity `json:"order"`
}

// Product defines model for Product.
type Product struct {
	Name string `json:"name"`
}

// ProductQuantity defines model for ProductQuantity.
type ProductQuantity struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

// Stock defines model for Stock.
type Stock struct {
	Stock []ProductQuantity `json:"stock"`
}

// Warehouse defines model for Warehouse.
type Warehouse struct {
	Id int `json:"id"`
}

// PostWarehousesJSONRequestBody defines body for PostWarehouses for application/json ContentType.
type PostWarehousesJSONRequestBody = AddWarehouse

// PostWarehousesWarehouseIdOrdersJSONRequestBody defines body for PostWarehousesWarehouseIdOrders for application/json ContentType.
type PostWarehousesWarehouseIdOrdersJSONRequestBody = CreateOrder

// PostWarehousesWarehouseIdProductsJSONRequestBody defines body for PostWarehousesWarehouseIdProducts for application/json ContentType.
type PostWarehousesWarehouseIdProductsJSONRequestBody = AddProductToWarehouse
