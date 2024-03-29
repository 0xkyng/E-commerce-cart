package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codekyng/E-commerce-cart.git/database"
	"github.com/codekyng/E-commerce-cart.git/models"
	generate "github.com/codekyng/E-commerce-cart.git/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

// HashPassword hashes user password
func HashPassword(password string) string {
	// Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// VerifyPassword compares user password & given password
func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	// Unhash the password and compare
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "login or password is incorrect"
		valid = false
	}

	return valid, msg

}

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Create user
		var user models.User
		// using BindJson method to serialize todo or extract data
		// From database to user
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the user struct
		validationErr := Validate.Struct(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		// Check if user email exist on database
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Checck count
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user already exist"})
		}

		// Check if user phone number exist on database
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Check count
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone number is already in use"})
			return
		}

		// User object
		// Password
		password := HashPassword(*user.Password)
		user.Password = &password

		// Timestamp
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// User Id
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		// Token
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken

		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		// Insert the user object above in the user collection.
		_, inserterr := UserCollection.InsertOne(ctx, user)
		// Handle error
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		defer cancel()

		// Return user
		c.JSON(http.StatusCreated, "Successfully signed in!👍")
	}
}

// Login Endpoint
func Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		// Verify if password exists
		passwordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		defer cancel()

		// Check if password is correct
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		// Generate token if user details are correct
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		// Update founduser details
		generate.UpdateAllToken(token, refreshToken, founduser.User_ID)

		// Return founduser
		c.JSON(http.StatusFound, founduser)

	}

}

// ProductViewerAdmin 
func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"not inserted"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfylly added")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create product list
		var productlist []models.Product
		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Find executes a find command and returns
		// A Cursor over the matching documents in the collection.
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "somrthing went wrong, please try after some time")
			return
		}

		// All iterates the cursor and decodes each document into results.
		// The results parameter must be a pointer to a slice.
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer cancel()

		// Return product list
		c.IndentedJSON(200, productlist)
	}
}

// SearchProductByQuery searches for a product by name
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var Searchproducts []models.Product
		queryParam := c.Query("name")

		// Check if name is empty
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		// set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})

		if err != nil {
			c.IndentedJSON(404, "something wentwrong while fetching the data")
			return
		}

		err = searchquerydb.All(ctx, &Searchproducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid request")
			return
		}

		defer cancel()

		// Return Searched product
		c.IndentedJSON(200, Searchproducts)
	}

}
