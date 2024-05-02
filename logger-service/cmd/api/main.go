package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WebPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://localhost:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	//connect to mongo

	mongoClient, err := ConnectToMongo()

	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	log.Println("starting service on port " + WebPort + "...")
	app := Config{
		Models: data.New(client),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WebPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

func ConnectToMongo() (*mongo.Client, error) {
	// create a connection options

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println(" Error connecting ", err)
		return nil, err
	}
	log.Println("Connected to MongoDB")
	return c, nil
}
