package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var products = []string{"Ciasto", "Placek", "Tort"}

const productNotFoundMsg = "Product not found"

func createProduct(c echo.Context) error {
	product := c.FormValue("product")
	products = append(products, product)
	return c.JSON(http.StatusCreated, product)
}

func getProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	product := products[productId]
	if product == "" {
		return c.JSON(http.StatusNotFound, productNotFoundMsg)
	}
	return c.JSON(http.StatusOK, product)
}

func updateProduct(c echo.Context) error {
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	product := c.FormValue("product")
	if products[productId] == "" {
		return c.JSON(http.StatusNotFound, productNotFoundMsg)
	}
	products[productId] = product
	return c.JSON(http.StatusOK, product)
}

func deleteProduct(c echo.Context) error {
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if products[productId] == "" {
		return c.JSON(http.StatusNotFound, productNotFoundMsg)
	}
	products = append(products[:productId], products[productId+1:]...)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	const productsIdUrl = "/products/:id"
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/products", createProduct)
	e.GET("/products", getProducts)
	e.GET(productsIdUrl, getProduct)
	e.PUT(productsIdUrl, updateProduct)
	e.DELETE(productsIdUrl, deleteProduct)

	e.Logger.Fatal(e.Start(":1323"))
}
