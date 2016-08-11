package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	idMap := map[string]map[string]int{}
	initialValue := 42
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
				fmt.Printf("Name %s was not found in Environment %s, adding it\n", name, environment)
				idMap[environment][name] = initialValue
				context.String(http.StatusOK, strconv.Itoa(initialValue))
			}
		} else {
			// add unfound environment and name
			fmt.Printf("Environment %s was not found, adding it\n", environment)
			idMap[environment] = map[string]int{name: initialValue}
			context.String(http.StatusOK, strconv.Itoa(initialValue))
		}
	})

	router.Run("localhost:8080")
}
