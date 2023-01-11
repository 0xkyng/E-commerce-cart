package tokens

import (
	"time"
	"log"
	"os"

	"github.com/codekyng/E-commerce-cart.git/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("CODE_KYNG")



func TokenGenerator(email string, firstname string, lastname string, uid string)(signedtoken string, signedrefreshtoken string, err error){

	claims := &SignedDetails{
		Email: email,
		First_Name: firstname,
		Last_Name: lastname,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	// Create token
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([](SECRET_KEY))
	
	if err != nil {
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshlaims).SignedString([](SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshtoken, err

}


func ValidateToken(){

}

func UpdateAllToken(){

}