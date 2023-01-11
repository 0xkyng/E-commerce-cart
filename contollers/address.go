package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/codekyng/E-commerce-cart.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"Invalid id"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addresses models.Address
		addresses.Address_id =primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Aggregation Pipeline Stages
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: 
		bson.D{primitive.E{Key: "sum", Value: 1}}}}}}

		// Run aggregate function
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal server Error")
		}

		// Create adddress info
		var addressinfo []bson.M
		err = pointcursor.All(ctx, &addressinfo)
		if err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}

		// Compare address sizes
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			c.IndentedJSON(400, "Not allowed")
		}
		defer cancel()
		ctx.Done()
	}

}


// EditHomeAddress edits user home address
func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query the user id whose
		// Home address is to be edited
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("Content-Type", "Application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server Error")
		}

		var editaddress models.Address
		if err = c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresss.0.house_name", Value: editaddress.House}, {Key: 
		"address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", 
		Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successsfully edited home address")
	}

}


// EditWorkAddress edist user work address
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query user id whose
		// Work address to be edited
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, "Invalid")
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address
		if err = c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresss.1.house_name", Value: editaddress.House}, {Key: 
		"address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: 
		"address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter,update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully edited work address")
	}

}

// DeleteAddress sets the address to be 
// Deleted to an empty value
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find the user id whose
		// You want to delete
		user_id := c.Query("id")

		// Check if user id is empty.
		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid Search Index"})
			c.Abort()
			return
		}

		// Create an empty slice of type address
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		// Set context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id", Value: usert_id}}
		// Set the user address to empty value
		update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Succesfylly Deleted")


	}

}