package main

import (
	dbConnection "Users/raghav.d/Desktop/jwtPoc/DBConnection"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	//"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type UserInputStruct struct {
	UserInput string `json:"UserInput" bson:"userInput"`
	Password  string `json:"Password" bson:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserName string `json:"userInput" bson:"userInput"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
}

type ResponseStruct struct {
	Token    string
	Name     string
	UserName string
	Email    string
}

type TokenResponse struct {
	Bearer string `json:"Bearer"`
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.POST("/signIn", postHandler)

	// Start server
	e.Logger.Fatal(e.Start(":81"))
}

func postHandler(req echo.Context) error {
	Request := UserInputStruct{}
	err := json.NewDecoder(req.Request().Body).Decode(&Request)
	if err != nil {
		return req.String(500, "error")
	}

	//res,err:= VerifyUserAndReturnResponse(Request.UserInput, Request.Password)
	if err != nil {
		return req.JSON(401, err)
	}

	tokenString, err := CreateToken(Request.UserInput, Request.Password)
	if err != nil {
		return req.JSON(401, err)
	}
	return req.JSON(200, TokenResponse{tokenString})
}

func CreateToken(userInput string, password string) (string, error) {
	var jwtKey = []byte("my_secret_key")
	expirationTime := time.Now().Add(15 + time.Minute)

	err := VerifyUserAndReturnResponse(userInput, password)

	if err != nil {
		return "", err
	}

	claims := &Claims{
		Username: userInput,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func VerifyUserAndReturnResponse(userInput string, password string) error {
	//making mongo connection
	mongoInstance := dbConnection.GetPool()
	rows := mongoInstance.Database("UserDB").Collection("UsersInfo")
	var elem UserInfo
	query1 := bson.M{"email": userInput}

	err := rows.FindOne(context.TODO(), query1).Decode(&elem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
		}
	}

	if elem.Password != password {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}
	return nil
}
