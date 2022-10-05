package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const uri = "mongodb+srv://root:1234@cluster0.ik76ncs.mongodb.net/?retryWrites=true&w=majority"

type user struct {
	ID    string `json:"id"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Password string `json:"password"`
	Age  int `json:"age"`
}

var users = []user{
	{ID: "1", Name: "John", Surname: "Doe", Email: "", Password: "dev.dilshodjon@gmail.com", Age: 20},
}


func createUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	//create users
	collection := client.Database("student").Collection("users")
	//insert user into collection users
	insertResult, err := collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
//get users from db and return users as json
func getUsers(c *gin.Context) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	//create users
	collection := client.Database("student").Collection("users")
	//insert user into collection users
	var users []user
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	for cur.Next(context.TODO()) {
		var elem user
		err := cur.Decode(&elem)
		if err != nil {
			panic(err)
		}
		users = append(users, elem)
	}
	if err := cur.Err(); err != nil {
		panic(err)
	}
	cur.Close(context.TODO())
	c.IndentedJSON(http.StatusOK, users)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser user
	if err := c.BindJSON(&updatedUser); err != nil {
		return
	}
	for i, a := range users {
		if a.ID == id {
			users[i] = updatedUser
		}
	}
	c.IndentedJSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, a := range users {
		if a.ID == id {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	c.IndentedJSON(http.StatusOK, users)
}

func main() {

	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.Run(":8080")
}
