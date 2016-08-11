package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	idMap := map[string]map[string]int{}
	incrementBy := 5

	router := gin.Default()

	router.GET("/:environment/:name", func(context *gin.Context) {
		name := context.Param("name")
		environment := context.Param("environment")
		// check if the environment is found, otherwise add it
		if _, ok := idMap[environment]; ok {
			// check if the name is found, otherwise add it
			if id, ok := idMap[environment][name]; ok {
				// increment id and return it
				idMap[environment][name] = id + incrementBy
				context.String(http.StatusOK, fmt.Sprintf("%d", idMap[environment][name]))
			} else {
				// add unfound name
				fmt.Println("Name ", name, " was not found in Environment ", environment, ", adding it")
				idMap[environment][name] = 1
				context.String(http.StatusOK, "1")
			}
		} else {
			// add unfound environment
			fmt.Println("Environment ", environment, " not found, adding it")
			idMap[environment] = map[string]int{name: 1}
			context.String(http.StatusOK, "1")
		}
	})

	router.Run("localhost:8080")
}
