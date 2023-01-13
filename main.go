package main

import (
	"log"
	"os"

	"github.com/codekyng/E-commerce-cart.git/controllers"
	"github.com/codekyng/E-commerce-cart.git/database"
	"github.com/codekyng/E-commerce-cart.git/middleware"
	"github.com/codekyng/E-commerce-cart.git/routes"
	"github.com/gin-gonic/gin"
)

// Application setup
func main() {
	port := os.Getenv("PORT")
	if port == ""{
		port ="8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveCartItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}