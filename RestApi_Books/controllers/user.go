package controllers

import (
	"errors"

	"modules/database"
	"modules/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type books_db struct {
	Db *gorm.DB
}

func New() *books_db {
	db := database.InitDb()
	db.AutoMigrate(&models.Book{})
	return &books_db{Db: db}
}

// create user
func (Books_db *books_db) CreateUser(c *gin.Context) {
	var book models.Book
	c.BindJSON(&book)
	err := models.CreateUser(Books_db.Db, &book)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, book)
}

// get users
func (Books_db *books_db) GetUsers(c *gin.Context) {
	var book []models.Book
	err := models.GetUsers(Books_db.Db, &book)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, book)
}

// get user by id
func (Books_db *books_db) GetUser(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	err := models.GetUser(Books_db.Db, &book, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, book)
}

// update user
func (Books_db *books_db) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	err := models.GetUser(Books_db.Db, &book, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&book)
	err = models.UpdateUser(Books_db.Db, &book)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, book)
}

// delete user
func (Books_db *books_db) DeleteUser(c *gin.Context) {
	var book models.Book
	id := c.Param("id")
	err := models.DeleteUser(Books_db.Db, &book, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
