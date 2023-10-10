package internal

import "fmt"

var (
	// ProductBook is a book.
	ProductBook = Product{"Book"}
	// ProductPen is a pen.
	ProductPen = Product{"Pen"}
)

// Set of all possible Product that can be ordered.
var products = map[string]Product{
	ProductBook.name: ProductBook,
	ProductPen.name:  ProductPen,
}

// Product represents a product that can be ordered.
// It is a value object and immutable.
type Product struct {
	name string
}

// Name returns the name of the product, never empty.
func (p Product) Name() string {
	return p.name
}

// FindProductByName finds a product by its name.
// It returns an error if the product is invalid.
func FindProductByName(value string) (Product, error) {
	product, exists := products[value]
	if !exists {
		return Product{}, fmt.Errorf("invalid product %q", value)
	}

	return product, nil
}
