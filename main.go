package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/:environment/:name", func(c *gin.Context) {
		name := c.Param("environment")
		environment := c.Param("name")
		c.String(http.StatusOK, "Hello %s in %s", name, environment)
	})

	router.Run("localhost:8080")
}
