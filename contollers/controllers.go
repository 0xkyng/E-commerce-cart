package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codekyng/E-commerce-cart.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

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
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		// Validate the user struct
		validationErr := Validate.Struct(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		// Check if user email exist on database
		count, err := UserCollection.countDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Checck count
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user already exist"})
		}

		// Check if user phone number exist on database
		count, err = UserCollection.countDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Check count
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"this phone number is already in use"})
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
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *&user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken

		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		// Insert the user object above in the user collection.
		_, inserterr := UserCollection.InsertOne(ctx, user)
		// Handle error
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"the user did not get created"})
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

		// Create user
		var user models.User
		// using BindJson method to serialize todo or extract data
       // From User struct to user
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err})
		}

		// Check if user email exists in the database
		err := UserCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&founduser)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
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
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.first_Name, *founduser.last_Name, *founduser.User_ID)

		// Update founduser details
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		// Return founduser
		c.JSON(http.StatusFound, founduser)


	}

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}

func SearchProductByQuery() {
	
}