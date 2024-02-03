package main

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func createTokens(guid string) (string, string) {
	accessToken, err := createAccessToken(guid)
	if err != nil {
		panic(err)
	}

	refreshToken, err := createRefreshToken(guid)
	if err != nil {
		panic(err)
	}

	return accessToken, refreshToken
}

func createAccessToken(guid string) (string, error) {
	expirationTimeAccessToken := time.Now().Add(15 * time.Minute).Unix()

	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTimeAccessToken
	claims["guid"] = guid

	if err := godotenv.Load("jwtService.env"); err != nil {
		log.Println("No .env file found")
	}
	key := []byte(os.Getenv("SECRET_ACCESS_KEY"))
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createRefreshToken(guid string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS512)

	expirationTimeRefreshToken := time.Now().Add(60 * time.Minute).Unix()

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = guid
	rtClaims["exp"] = expirationTimeRefreshToken

	if err := godotenv.Load("jwtService.env"); err != nil {
		log.Println("No .env file found")
	}
	key := []byte(os.Getenv("SECRET_REFRESH_KEY"))
	refreshTokenString, err := refreshToken.SignedString(key)
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}
