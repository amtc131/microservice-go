package main

import (
	"fmt"
	"microservice/main.go/client/client"
	"microservice/main.go/client/client/products"
	"testing"
)

func TestOutClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod)
	t.Fail()
}
