package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_NAME            = "instagram"
	DB_USER_COLLECTION = "user"
	DB_POST_COLLECTION = "post"
)

type User struct {
	Id       primitive.ObjectID `json:"Id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"Name" bson:"Name"`
	Email    string             `json:"Email" bson:"Email"`
	Password string             `json:"Password" bson:"Password"`
}

type Post struct {
	Id               primitive.ObjectID  `json:"Id,omitempty" bson:"_id,omitempty"`
	Uid              primitive.ObjectID  `json:"UId,omitempty" bson:"Uid,omitempty"`
	Caption          string              `json:"Caption,omitempty"`
	Image_URL        string              `json:"Image_URL,omitempty"`
	Posted_Timestamp primitive.Timestamp `json:"Timestamp,omitempty"`
}

var client *mongo.Client

func reportError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message":"` + err.Error() + `"}`))

}

func getIDParam(url string) string {
	p := strings.Split(url, "/")
	return p[len(p)-1]
}

func userGETHandler(response http.ResponseWriter, request *http.Request) {
	// Set the respone content type to JSON
	response.Header().Set("Content-Type", "application/json")
	//The collection to store new users
	collection := client.Database(DB_NAME).Collection(DB_USER_COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		// Extract id parameter from URL
		uid := getIDParam(request.URL.Path)
		// Converting string id to ObjectID
		id, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			reportError(response, err)
		}
		var user User
		// Filter to search from user with id.
		idFilter := bson.M{"_id": id}
		err = collection.FindOne(ctx, idFilter).Decode(&user)
		if err != nil {
			// If user doesnt exists
			reportError(response, errors.New("User Doesnt Exists"))
			return
		}
		json.NewEncoder(response).Encode(user)

	}
func userPOSTHandler(response http.ResponseWriter,request *http.Request){
	// Set the respone content type to JSON
	response.Header().Set("Content-Type", "application/json")
	//The collection to store new users
	collection := client.Database(DB_NAME).Collection(DB_USER_COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
var user User
		// Convert the json payload into go understandable data structure
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			reportError(response, err)
			return
		}
		//Searching if the user with user name already existes
		emailFilter := bson.M{"Email": user.Email}
		fCursor, err := collection.Find(ctx, emailFilter)
		if err != nil {
			reportError(response, err)
		}
		// slice to store the user with the same email,in any
		var filteredUser []bson.M
		if err = fCursor.All(ctx, &filteredUser); err != nil {
			reportError(response, err)
			return
		}
		defer fCursor.Close(ctx)
		// if user with same email exists
		if len(filteredUser) != 0 {
			var alreadyExistsError = errors.New("User Already Exists")
			reportError(response, alreadyExistsError)
			return
		}
		hash := sha256.New()
		hash.Write([]byte(user.Password))
		user.Password = hex.EncodeToString(hash.Sum(nil))
		result, _ := collection.InsertOne(ctx, user)
		json.NewEncoder(response).Encode(result)
}
		


func postGETHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	pcollection := client.Database(DB_NAME).Collection(DB_POST_COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		//Checking for id in query parameters
		pid := getIDParam(request.URL.Path)
		id, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			reportError(response, err)
		}
		var post Post
		// Filter the post to find post with 'id'
		idFilter := bson.M{"_id": id}
		err = pcollection.FindOne(ctx, idFilter).Decode(&post)
		if err != nil {
			reportError(response, errors.New("Post Doesnt Exists"))
			return
		}
		json.NewEncoder(response).Encode(post)
	}
func postPOSTHandler(response http.ResponseWriter, request *http.Request){
	ucollection := client.Database(DB_NAME).Collection(DB_USER_COLLECTION)
	pcollection := client.Database(DB_NAME).Collection(DB_POST_COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		// Only users with an account can make a post.
		var post Post
		err := json.NewDecoder(request.Body).Decode(&post)
		if err != nil {
			reportError(response, err)
			return
		}
		//Verify is the user with uid exists
		var user User
		idFilter := bson.M{"_id": post.Uid}
		// Searching for the user with 'Uid'
		err = ucollection.FindOne(ctx, idFilter).Decode(&user)
		if err != nil {
			// If user doesn't exists, raise error.
			reportError(response, errors.New("User Doesnt Exists"))
			return
		}
		// Store the unix timestamp(milliseconds from epoch)
		post.Posted_Timestamp = primitive.Timestamp{T: uint32(time.Now().Unix())}
		result, err := pcollection.InsertOne(ctx, post)
		if err != nil {
			reportError(response, err)
			return
		}
		json.NewEncoder(response).Encode(result)
}

func userPostPOSTHandler(response http.ResponseWriter, request *http.Request) {
	log.Println(request.URL)
	response.Header().Set("Content-Type", "application/json")
	collection := client.Database(DB_NAME).Collection(DB_POST_COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// End point only handles GET, check for get method
	if request.Method != "GET" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	uid := getIDParam(request.URL.Path)
	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		reportError(response, err)
		return
	}
	// Filter posts with Uid as 'id'
	postFilter := bson.M{"Uid": id}
	fCursor, err := collection.Find(ctx, postFilter)
	if err != nil {
		reportError(response, err)
		return
	}
	defer fCursor.Close(ctx)
	var userPosts []bson.M
	if err = fCursor.All(ctx, &userPosts); err != nil {
		reportError(response, err)
		return
	}
	//TODO Pagination of posts by user with UID 'id'
	json.NewEncoder(response).Encode(userPosts)
}

func main() {
	//Connection to the mongod server listining at 270217
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOpts := options.Client().ApplyURI(
		"mongodb://localhost:27017/?connect=direct")
	var err error
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Error :%v", err)
	}

	// URL Handlers
	http.HandleFunc("/users/", userGETHandler)
	http.HandleFunc("/users",userPOSTHandler)
	http.HandleFunc("/posts/", postGETHandler)
	http.HandleFunc("/posts",postPOSTHandler)
	http.HandleFunc("/posts/users/", userPostPOSTHandler)
	// Server to run at port 8080
	http.ListenAndServe(":8080", nil)
}
