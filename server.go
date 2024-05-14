package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var products = []string{"Ciasto", "Placek", "Tort"}

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
		return c.JSON(http.StatusNotFound, "Product not found")
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
		return c.JSON(http.StatusNotFound, "Product not found")
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
		return c.JSON(http.StatusNotFound, "Product not found")
	}
	products = append(products[:productId], products[productId+1:]...)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/products", createProduct)
	e.GET("/products", getProducts)
	e.GET("/products/:id", getProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)

	e.Logger.Fatal(e.Start(":1323"))
}
