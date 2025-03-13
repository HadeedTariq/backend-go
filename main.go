package main

import (
	"context"
	"fmt"
	"log"
	"my-backend/config"
	"my-backend/controller"
	"my-backend/db"
	"my-backend/middlewares"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	env, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading envs:", err)
	}

	client, err := db.ConnectToDb(env.DBURI)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("mongoClient", client)
		c.Next()
	})

	router.Use(middlewares.LoggingMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	router.GET("/users", func(c *gin.Context) {
		type User struct {
			ID    string `json:"id" bson:"_id,omitempty"`
			Name  string `json:"name" bson:"name"`
			Email string `json:"email" bson:"email"`
		}

		client, exists := c.Get("mongoClient")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
			return
		}

		mongoClient, ok := client.(*mongo.Client)

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
			return
		}

		collection := mongoClient.Database("courseGarden").Collection("users")

		cursor, err := collection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't find collection"})
			return
		}

		defer cursor.Close(context.TODO())

		var users []User

		for cursor.Next(context.TODO()) {
			var user User
			if err := cursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Collections are not correct"})
				return
			}
			users = append(users, user)
		}

		c.JSON(200, gin.H{"message": "MongoDB connection works!", "collection": collection.Name(), "users": users})
	})

	router.POST("/auth/register", controller.RegisterUser)

	go func() {
		if err := router.Run(":3000"); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	fmt.Println("Shutting down server...")

	// Release resources
	db.DisconnectMongoDB()
	fmt.Println("Server gracefully stopped.")

}
