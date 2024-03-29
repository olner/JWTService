package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	service "jwtService/internal/services"
)

var connection *mongo.Client

func CreateConncetion() {
	if err := godotenv.Load("jwtService.env"); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	connection = client
	log.Println("Connection open")
}

func CloseConnection() {
	if err := connection.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	log.Println("Connection closed")
}

func IsCollectionExists() bool {
	cNames, err := connection.Database("test").ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, value := range cNames {
		if value == "name" {
			return true
		}
	}
	return false
}

func CreateCollection() {
	db := connection.Database("test")
	err := db.CreateCollection(context.TODO(), "tokens")
	if err != nil {
		log.Fatal(err)
	}
}

func SaveRefreshToken(guid string, token string) {
	CreateConncetion()
	defer func() {
		CloseConnection()
	}()

	if !IsCollectionExists() {
		CreateCollection()
	}

	coll := connection.Database("test").Collection("tokens")
	var result bson.M
	coll.FindOne(context.TODO(), bson.D{{Key: "guid", Value: guid}}).Decode(&result)

	hashOfRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	if len(result) == 0 {
		coll.InsertOne(
			context.TODO(),
			bson.D{
				{Key: "guid", Value: guid},
				{Key: "refreshToken", Value: string(hashOfRefreshToken)},
			},
		)
		log.Println("The refresh token is saved in the DB.")
	} else {
		log.Println("Duplicate! This guid parameter is already registered")
	}
}

func UpdateRefreshToken(guid string, token string) {
	CreateConncetion()
	defer func() {
		CloseConnection()
	}()

	newRefreshToken, err := service.CreateRefreshToken(guid)
	if err != nil {
		panic(err)
	}
	newHashOfRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	coll := connection.Database("test").Collection("tokens")
	_, updateErr := coll.UpdateOne(
		context.TODO(),
		bson.D{{Key: "guid", Value: guid}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "refreshToken", Value: string(newHashOfRefreshToken)}}}},
	)
	if updateErr != nil {
		log.Println("Update error")
	}
}

func IsValidateRefreshToken(guid string, refreshToken string) bool {
	CreateConncetion()
	defer func() {
		CloseConnection()
	}()

	coll := connection.Database("test").Collection("tokens")
	var hashOfRefreshTokenFromDB bson.M
	errFromDB := coll.FindOne(context.TODO(), bson.D{{Key: "guid", Value: guid}}).Decode(&hashOfRefreshTokenFromDB)

	if errFromDB != nil {
		log.Println("Refresh tokens did not match")
		return false
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(hashOfRefreshTokenFromDB["refreshToken"].(string)), []byte(refreshToken))
		if err != nil {
			log.Println("Refresh tokens did not match")
			return false
		}
	}
	return true
}
