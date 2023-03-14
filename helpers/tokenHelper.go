package helpers

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"jwt-with-go/database"
	"os"
	"time"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.StandardClaims
}

var userCollection = database.OpenCollection(database.Client, "users")

func GenerateAllTokens(email string, firstname string, lastname string, usertype string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Uid:       uid,
		UserType:  usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		Email:     "",
		FirstName: "",
		LastName:  "",
		Uid:       "",
		UserType:  "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY ")))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(SignedDetails)
	if !ok {
		msg = fmt.Sprintf("The token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
	}

	return claims, msg

}

func UpdateAllTokens(signedToken, signedRefreshToken, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})
	UpdatedAt, _ := time.Parse(time.RFC850, time.Now().Format(time.RFC850))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: UpdatedAt})

	upsert := true
	filter := bson.M{"userId": userId}
	filterOptions := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &filterOptions)

	if err != nil {
		fmt.Println(err.Error())
	}
	return
}
