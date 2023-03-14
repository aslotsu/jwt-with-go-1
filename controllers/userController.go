package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"jwt-with-go/database"
	"jwt-with-go/helpers"
	"jwt-with-go/models"
	"time"
)

var userCollection = database.OpenCollection(database.Client, "reg-users")

var validate = validator.New()

func getUserHandlerFunc(c *gin.Context) {
	userId := c.Param("user_id")
	if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User

	if err := userCollection.FindOne(ctx, bson.M{"userId": userId}).Decode(user); err != nil {
		err = errors.New("cannot find user")
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		err = errors.New("could not hash the password, password is unsafe to use ")
	}

	return string(bytes)

}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check, msg := true, ""
	if err != nil {
		check = false
		msg = "password does not match this account"
	}
	return check, msg

}

func signUpHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	defer cancel()

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if validationErr := validate.Struct(user); validationErr != nil {
		fmt.Println(validationErr.Error())
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		fmt.Println(err.Error())
	}
	if count > 0 {
		c.JSON(409, gin.H{"err": "This email already exists in the database"})
		return
	}
	password := HashPassword(user.Password)
	user.Password = password

	phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phonenumber": user.PhoneNumber})
	defer cancel()
	if err != nil {
		fmt.Println(err.Error())
	}
	if phoneCount > 0 {
		c.JSON(409, gin.H{"err": "This phone number exists in the database"})
		return
	}

	user.CreatedAt, err = time.Parse(time.RFC850, time.Now().Format(time.RFC850))
	if err != nil {
		fmt.Println(err.Error())
	}

	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.UserType, user.UserId)
	user.Token = token
	user.RefreshToken = refreshToken
	insertedId, err := userCollection.InsertOne(ctx, user)
	c.JSON(201, insertedId)
}

func loginHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	var foundUser models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}
	passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
	if !passwordIsValid {
		c.JSON(303, gin.H{"err": msg})
		return
	}

	if foundUser.Email == "" {
		c.JSON(401, gin.H{"err": "Could not find user with this email"})
		return
	}
	token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserType, foundUser.UserId)
	foundUser.UpdatedAt, _ = time.Parse(time.RFC850, time.Now().Format(time.RFC850))
	helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)

	if err := userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserId}).Decode(&foundUser); err != nil {
		c.JSON(405, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, foundUser)

}

func getAllUsersHandlerFunc(c *gin.Context) {
	if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		c.JSON(403, gin.H{"err": err.Error()})
		return
	}

}

func SignUp() gin.HandlerFunc {
	return signUpHandlerFunc
}

func Login() gin.HandlerFunc {
	return loginHandlerFunc
}

func GetUsers() gin.HandlerFunc {
	return getAllUsersHandlerFunc
}

func GetUser() gin.HandlerFunc {
	return getUserHandlerFunc
}
