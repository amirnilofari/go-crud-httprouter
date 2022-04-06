package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var portNumber = ":9000"
var mongoURI = "mongodb://localhost:27017"
var client *mongo.Client

type Post struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title,omitempty"`
	Body  string             `json:"body" bson:"body,omitempty"`
}

func CreatePost(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.Header().Add("Content-Type", "application/json")
	var post Post

	json.NewDecoder(request.Body).Decode(&post)

	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, _ := collection.InsertOne(ctx, post)

	json.NewEncoder(response).Encode(result)
}

func GetPosts(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")
	var posts []Post
	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, currErr := collection.Find(ctx, bson.D{})
	if currErr != nil {
		panic(currErr)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &posts); err != nil {
		panic(err)
	}

	json.NewEncoder(response).Encode(posts)

}

func GetPost(response http.ResponseWriter, request *http.Request, param httprouter.Params) {
	response.Header().Add("Content-Type", "application/json")

	id, _ := primitive.ObjectIDFromHex(param.ByName("id"))

	var post Post

	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := collection.FindOne(ctx, Post{
		ID: id,
	}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(post)

}

func DeletePost(response http.ResponseWriter, request *http.Request, param httprouter.Params) {
	response.Header().Add("Content-Type", "application/json")

	id, _ := primitive.ObjectIDFromHex(param.ByName("id"))

	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := collection.DeleteOne(ctx, Post{ID: id})
	if err != nil {
		panic(err)
	}
	json.NewEncoder(response).Encode(result)

}

func UpdatePost(response http.ResponseWriter, request *http.Request, param httprouter.Params) {
	response.Header().Add("Content-Type", "application/json")

	id, _ := primitive.ObjectIDFromHex(param.ByName("id"))
	var post Post
	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	json.NewEncoder(response).Encode(post)

}

func main() {
	fmt.Println("Starting the application")
	//Connect to MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	server := httprouter.New()

	server.GET("/posts", GetPosts)
	server.GET("/posts/:id", GetPost)

	server.POST("/post", CreatePost)

	server.DELETE("/posts/:id", DeletePost)

	server.PUT("/posts/:id", UpdatePost)

	log.Fatal(http.ListenAndServe(portNumber, server))
}
