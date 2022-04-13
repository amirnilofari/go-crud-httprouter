package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amirnilofari/crud-httprouter-mongo/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var client *mongo.Client

func CreatePost(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.Header().Add("Content-Type", "application/json")
	var post models.Post

	json.NewDecoder(request.Body).Decode(&post)

	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, _ := collection.InsertOne(ctx, post)

	json.NewEncoder(response).Encode(result)
}

func GetPosts(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	response.Header().Set("Content-Type", "application/json")
	var posts []models.Post
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

	var post models.Post

	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := collection.FindOne(ctx, models.Post{
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

	result, err := collection.DeleteOne(ctx, models.Post{ID: id})
	if err != nil {
		panic(err)
	}
	json.NewEncoder(response).Encode(result)

}

func UpdatePost(response http.ResponseWriter, request *http.Request, param httprouter.Params) {
	response.Header().Add("Content-Type", "application/json")

	id, _ := primitive.ObjectIDFromHex(param.ByName("id"))

	var post models.Post
	json.NewDecoder(request.Body).Decode(&post)

	collection := client.Database("blog").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"title", post.Title}, {"body", post.Body}}},
		})
	if err != nil {
		fmt.Println("ERROR IS : ", err)
	}

	json.NewEncoder(response).Encode(result)

}
