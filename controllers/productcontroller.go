package controllers

import (
	"backendGO/database"
	"backendGO/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, price, user_id FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": products})
}

func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	err = database.DB.QueryRow("SELECT id, name, price, user_id FROM products WHERE id=$1", id).Scan(&product.ID, &product.Name, &product.Price, &product.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var newProduct models.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO products (name, price, user_id) VALUES ($1, $2, $3) RETURNING id",
		newProduct.Name, newProduct.Price, userID,
	).Scan(&newProduct.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newProduct})
}

func UpdateProduct(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var productOwnerID string
	err = database.DB.QueryRow("SELECT user_id FROM products WHERE id=$1", id).Scan(&productOwnerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if productOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	_, err = database.DB.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3 AND user_id=$4", updatedProduct.Name, updatedProduct.Price, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedProduct})
}

func DeleteProduct(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var productOwnerID string
	err = database.DB.QueryRow("SELECT user_id FROM products WHERE id=$1", id).Scan(&productOwnerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if productOwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM products WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
