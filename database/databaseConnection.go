package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"go.mongo.org/mongo-driver/mongo"

)
func DBinstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Eroor loading .env file")
	}

	
}
