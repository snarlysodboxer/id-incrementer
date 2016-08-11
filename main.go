package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var idMap = map[string]map[string]int{}
var initialValue = 42
var incrementBy = 5

func getOrSetID(name, environment string, passedID int) (int, string) {
	// check if the environment is found, otherwise add it
	if _, ok := idMap[environment]; ok {
		// check if the name is found, otherwise add it
		if foundID, ok := idMap[environment][name]; ok {
			if passedID == 0 {
				idMap[environment][name] = foundID + incrementBy
			} else {
				idMap[environment][name] = passedID
			}
		} else {
			// add unfound name, return ID
			fmt.Printf("Name %s was not found in Environment %s, adding it\n", name, environment)
			if passedID == 0 {
				idMap[environment][name] = initialValue
			} else {
				idMap[environment][name] = passedID
			}
		}
	} else {
		// add unfound environment and name, return ID
		fmt.Printf("Environment %s was not found, adding it\n", environment)
		if passedID == 0 {
			idMap[environment] = map[string]int{name: initialValue}
		} else {
			idMap[environment] = map[string]int{name: passedID}
		}
	}
	return http.StatusOK, strconv.Itoa(idMap[environment][name])
}

func main() {
	router := gin.Default()

	router.POST("/setter", func(context *gin.Context) {
		passedID, err := strconv.Atoi(context.PostForm("id"))
		if err != nil {
			context.String(http.StatusBadRequest, fmt.Sprintf("Error converting %s to an integer", context.PostForm("id")))
			return
		}
		if passedID == 0 {
			context.String(http.StatusBadRequest, "You can't set the ID to zero") // TODO rethink this limitation?
		} else {
			status, id := getOrSetID(context.PostForm("name"), context.PostForm("environment"), passedID)
			context.String(status, id)
		}
	})

	router.GET("/:environment/:name", func(context *gin.Context) {
		status, id := getOrSetID(context.Param("name"), context.Param("environment"), 0)
		context.String(status, id)
	})

	router.Run("localhost:8080")
}
