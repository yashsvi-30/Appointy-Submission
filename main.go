package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	//"strconv"
	//"log"
	"github.com/gorilla/mux"
	//"github.com/bmizerany/pat"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var Serve http.Handler
var client *mongo.Client

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name 	  string             `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumber    string		 `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	DateofBirth  string          `json:"dateofbirth,omitempty" bson:"dateofbirth,omitempty"`
	EmailAddress   string		 `json:"emailaddress,omitempty" bson:"emailaddress,omitempty"`
 	CreationTimestamp   *time.Time `json:"creation_timestamp,omitempty" bson:"creationtimestamp,omitempty"`
}

type Contact struct {
	IDone  	primitive.ObjectID  `json:"_idone,omitempty" bson:"_idone,omitempty"`
	IDtwo  	primitive.ObjectID  `json:"_idtwo,omitempty" bson:"_idtwo,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users Person
	_ = json.NewDecoder(request.Body).Decode(&users)
	collection := client.Database("yashsvisharma").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, users)
	json.NewEncoder(response).Encode(result)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var users Person
	collection := client.Database("yashsvisharma").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	
	//----------Tried using pat but its not working-------------------//

	//id, _ := primitive.ObjectIDFromHex(request.URL.Query().Get(":id")
	//objectIDS, _ := primitive.ObjectIDFromHex(string(id))
	//filter := bson.M{"_id": objectIDS}
	//err := collection.FindOne(ctx, filter).Decode(&users)
	//collection := client.Database("yashsvisharma").Collection("user")
	
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&users)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(users)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router :=mux.NewRouter()
	router.HandleFunc("/users", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/users/{id}", GetPersonEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)

	//--------------------Tried using pat but its not working---------------------------//

	//router := pat.New()
	//router.Post("/users", http.HandlerFunc(CreatePersonEndpoint))
	//router.Get("/users/{id}", http.HandlerFunc(GetPersonEndpoint))
	

	//--------------------Tried using httpRouter, but that too is not working---------------------//

	//http.HandleFunc("/",CreatePersonEndpoint)
	//http.HandleFunc("/",GetPersonEndpoint)
    //log.Fatal(http.ListenAndServe(":12345", nil))
}