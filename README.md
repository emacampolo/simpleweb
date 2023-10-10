# Simple web example

## Build the application

```bash
make # will build the application and run all tests.
```

If you only want to run the tests, you can run the following command:

```bash
make test
```

## Run the application

After cloning the application, cd into the root directory and run the following command:

```bash
make run
```

You can stop the application with `Ctrl+C`.

### cURLs for interacting with the application

After running, `make run`, open a new terminal and run the following commands:

```bash
# Add a new warehouse. 
curl -X POST -H "Content-Type: application/json" -s -d '{"name":"my warehouse"}' http://localhost:8080/warehouses

# Increment the quantity of a product in a warehouse.
curl -X POST -H "Content-Type: application/json" -s -d '{"product_name":"Book","quantity": 10}' http://localhost:8080/warehouses/1/products

# Create an order.
curl -X POST -H "Content-Type: application/json" -s -d '{"products":[{"product_name":"Book","quantity":1}]}' http://localhost:8080/warehouses/1/orders
```