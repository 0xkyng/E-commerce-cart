package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/codekyng/E-commerce-cart.git/database"
	"github.com/codekyng/E-commerce-cart.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		// Run the database level function

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully added to cart")

	}

}

// RemoveItem removes item from the cart
func (app *Application) RemoveCartItem() gin.HandlerFunc {
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

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		// Call the database level function
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully removed item from cart")
	}
}

// GetItemFromCart selects a particular item from the cart
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid id"})
			c.Abort()
			return
		}

		// Get user id
		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get cart details for th user
		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		// Aggregation Pipeline Stages
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: user_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: 
		"$sum", Value: "$usercat.price"}}}}}}

		// Run aggregation function
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		err = pointcursor.All(ctx, &listing)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		// Range over the data
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()

	}

}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user exists
		userQueryID := c.Query("id")

		if userQueryID == "" {
			log.Panicln("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Call the database level function
		err := database.BuyFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON("Successfully placed order")

	}

}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user id exists
		UserQueryID := c.Query("id")
		if UserQueryID == "" {
			log.Println("user id is empty")
			// Abort program
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		// Check if product id exists
		ProductQueryID := c.Query("id")
		if ProductQueryID == "" {
			log.Println("Product id id is empty")
			// Abort program
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		// //Check if the product is genuine
		productID, err := primitive.ObjectIDFromHex(ProductQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		// Call the database level function
		err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productID, UserQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successully placed the order")
	}
}
