package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Rule      string `json:"rule"`
	Password  string `json:"password"`
}

func main() {
	client := connectToMongo()
	createDatabse(client)
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.Password = HashPassword(user.Password)
		fmt.Println(user.Password)
		x := CheckPasswordHash("ali123", user.Password)
		createUser(client, user)
		fmt.Println(x)
		c.JSON(http.StatusOK, gin.H{"status": "added"})
	})

	r.GET("/print", func(c *gin.Context) {
		returnCollection(client, "users")
	})
	r.GET("/list", func(c *gin.Context) {
		returnAllCollections(client)
	})
	r.Run()
}

func returnCollection(client *mongo.Client, colName string) {
	collection := client.Database("testDB").Collection(colName)
	var results []bson.M
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func returnAllCollections(client *mongo.Client) {
	names, err := client.Database("testDB").ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	//var collection *mongo.Collection
	for _, x := range names {
		fmt.Println(x)
		collection := client.Database("testDB").Collection(x)
		var results []bson.M
		cur, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		if err = cur.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
		fmt.Println("------------")

	}
	//var results []bson.M
	//cur, err := collection.Find(context.TODO(), bson.M{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err = cur.All(context.TODO(), &results); err != nil {
	//	log.Fatal(err)
	//}
	//for _, result := range results {
	//	fmt.Println(result)
	//}
}

func createUser(client *mongo.Client, user User) {

	res, err := insertOne(client, context.Background(), "testDB", "users", user)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
