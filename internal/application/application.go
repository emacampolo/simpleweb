package application

import (
	"bootcamp-web/internal"
	"bootcamp-web/internal/handler"
	"bootcamp-web/internal/handler/mid"
	"bootcamp-web/internal/storage/storagememory"
	"bootcamp-web/platform/web"
	"context"
	"log"
	"net"
	"net/http"
	"os"
)

// Application represents the application running in a http server
// with all its dependencies.
type Application struct {
	server  *http.Server
	network string
	address string
}

// New creates a new un-started application.
func New() *Application {
	muxer := web.NewMux(mid.NewError(), mid.NewPanic())
	registerRoutes(storagememory.NewWarehouseRepository(), muxer)
	httpServer := &http.Server{Handler: muxer}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Application{
		server:  httpServer,
		network: "tcp",
		address: ":" + port,
	}
}

// Start starts the application and blocks until the application is stopped.
func (a *Application) Start() error {
	ln, err := net.Listen(a.network, a.address)
	if err != nil {
		return err
	}

	log.Printf("Started at %s\n", ln.Addr().String())

	return a.server.Serve(ln)
}

// Stop stops the application.
func (a *Application) Stop() error {
	return a.server.Shutdown(context.Background())
}

func registerRoutes(warehouseRepository internal.WarehouseRepository, m *web.Muxer) {
	m.Handle("POST", "/warehouses",
		handler.NewAddWarehouse(warehouseRepository))
	m.Handle("POST", "/warehouses/{warehouse_id}/products",
		handler.NewAddProductStock(warehouseRepository))
	m.Handle("POST", "/warehouses/{warehouse_id}/orders",
		handler.NewCreateOrder(warehouseRepository, internal.NewOrderService(warehouseRepository)))
}
