package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Begc007/url-short/core"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	router := gin.Default()
	router.POST("/url", addNewURL)
	router.GET("/url/:shortened", redirectURL)
	router.GET("/url/all", getAllURLs)
	router.Run(":8080")

}

func addNewURL(c *gin.Context) {
	var url core.URL
	err := c.ShouldBindJSON(&url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("URL: ", url)
	shorted := core.GenerateShortKey()
	log.Println("Shortened: ", shorted)

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Example: Set a key-value pair
	ctx := context.Background()
	err = client.Set(ctx, shorted, url.Value, 0).Err()
	if err != nil {
		log.Println("Error saving redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"shortened": shorted})
}

func getURL(c *gin.Context) {
	shortened := c.Param("shortened")
	if shortened == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shortened URL is required"})
		return
	}

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	ctx := context.Background()
	val, err := client.Get(ctx, shortened).Result()

	if err != nil {
		log.Println("Error getting redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"url": val})
}

func redirectURL(c *gin.Context) {
	shortened := c.Param("shortened")
	if shortened == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shortened URL is required"})
		return
	}

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	ctx := context.Background()
	val, err := client.Get(ctx, shortened).Result()

	if err != nil {
		log.Println("Error getting redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(val)
	c.Redirect(http.StatusFound, "http://"+val)
}

func getAllURLs(c *gin.Context) {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	ctx := context.Background()
	val, err := client.Keys(ctx, "*").Result()

	if err != nil {
		log.Println("Error getting redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"urls": val})
}
