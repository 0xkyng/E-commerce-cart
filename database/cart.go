package database

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrCantFindProduct = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find the product")
	ErrUserIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from cart")
	ErrCantBUyItemcart = errors.New("can't update the purchase")
)

func AddProductToCart() gin.HandlerFunc{

}

func RemoveCartItem() gin.HandlerFunc {

}

func GetItemFromCart() gin.HandlerFunc{

}

func BuyFromCart() gin.HandlerFunc {

}

func InstantBuy() gin.HandlerFunc {

}