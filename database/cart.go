package database

import (
	"errors"
	"log"
	"time"

	"github.com/codekyng/E-commerce-cart.git/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from cart")
	ErrCantBUyItemcart = errors.New("can't update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection, *mongo.Collection,  productID primitive.ObjectID, userID string) error{
	// Search for the product to be
	// To the cart from thr db
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Panic(err)
		return ErrCantFindProduct
	}

	var productcart []models.ProductUser
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	// Check user id
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIsNotValid
	}

	// Update user cart from db
	filter :=  bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", 
	Value: productCart}}}}}}

	_, err = UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil

}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	// Find item from Db
	// And remove from the cart
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil

}

func GetItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error{
	// Get user id
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems models.User
	var odercart models.Order

	// Checkout cart
	odercart.Order_ID = primitive.NewObjectID()
	odercart.Ordered_At = time.Now()
	odercart.Order_Cart = make([]models.ProductUser, 0)
	odercart.Payment_Method.COD = true

	// Add product user details
	// By running the aggregate function
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", 
	Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	currentresults, err := userCollection.Aggregate(ctx, unwind, grouping)
	ctx.Done()
	if err != nil {
		panic(err)
	}
	// Fetch the cart of the user
	// Find the cart tool
	// Create an order with the items
	// Empty the cart

}

func BuyFromCart() gin.HandlerFunc {

}

func InstantBuy() gin.HandlerFunc {

}