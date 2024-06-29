package controllers

import (
	"backendGO/database"
	"backendGO/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Cek apakah username atau email sudah terdaftar
	if err := database.DB.QueryRow("SELECT id FROM users WHERE username=$1 OR email=$2", newUser.Username, newUser.Email).Scan(&newUser.ID); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Hash password sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	newUser.PasswordHash = string(hashedPassword)

	// Simpan pengguna ke dalam database
	_, err = database.DB.Exec("INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)", newUser.Username, newUser.PasswordHash, newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newUser})
}

func Login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var storedUser models.User
	err := database.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username=$1", creds.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Bandingkan password yang di-input dengan password yang sudah di-hash
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token untuk autentikasi

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successful"})
}
