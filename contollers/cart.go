package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/codekyng/E-commerce-cart.git/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

// AddToCart adds product to tehe cart
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if product id exists
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			// Abort program
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		// Check if user id exists
		userQueryID := c.Query("UserID")
		if userQueryID == "" {
			log.Println("user id is empty")

			// Abort the program
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		// Check if product is genuine
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Run the database level function
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully added to cart")

	}

}

// RemoveItem removes item from the cart
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if product id exists
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			// Abort program
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		// Check if user id exists
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			// Abort the program
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		//Check if the product is genuine
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Call the database level function
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err =database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully removed item from cart")



	}

}

func GetItemFromCart() gin.HandlerFunc {

}

func BuyFromCart() gin.HandlerFunc {

}

func InstantBuy() gin.HandlerFunc {

}
