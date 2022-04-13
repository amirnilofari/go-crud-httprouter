package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amirnilofari/crud-httprouter-mongo/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var portNumber = ":9000"
var mongoURI = "mongodb://localhost:27017"

//var client *mongo.Client

func main() {
	fmt.Println("Starting the application")
	//Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	server := httprouter.New()

	server.GET("/posts", controllers.GetPosts)
	server.GET("/posts/:id", controllers.GetPost)

	server.POST("/post", controllers.CreatePost)

	server.DELETE("/posts/:id", controllers.DeletePost)

	server.PUT("/posts/:id", controllers.UpdatePost)

	log.Fatal(http.ListenAndServe(portNumber, server))
}
