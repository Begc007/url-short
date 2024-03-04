package main

import (
	"log"
	"net/http"

	"github.com/Begc007/url-short/core"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/shorten", addNewURL)
	router.Run(":8080")
	// url := "https://docs.github.com/en/copilot/using-github-copilot/getting-started-with-github-copilot?tool=vscode"
	// shorted := core.GetShortened(url)
	// fmt.Println(shorted)
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

	//insert into redis

	c.JSON(200, gin.H{"shortened": shorted})
}
