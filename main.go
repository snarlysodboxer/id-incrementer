package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

// TODO setup auto API documentation, add auth, add persistent storage, add settings flags or file

var idMap = map[string]map[string]int{}
var initialValue = 42
var incrementBy = 5
var mutex = &sync.Mutex{}

func getID(name, environment string) (int, string) {
	// check if the environment is found
	if _, ok := idMap[environment]; ok {
		// check if the name is found
		if id, ok := idMap[environment][name]; ok {
			idMap[environment][name] = id + incrementBy
		} else {
			// add unfound name
			fmt.Printf("Adding `%s/%s` with initial value `%d`\n", environment, name, initialValue)
			idMap[environment][name] = initialValue
		}
	} else {
		// add unfound environment and name
		fmt.Printf("Adding `%s/%s` with initial value `%d`\n", environment, name, initialValue)
		idMap[environment] = map[string]int{name: initialValue}
	}
	return http.StatusOK, strconv.Itoa(idMap[environment][name])
}

func setID(name, environment string, id int) (int, string) {
	fmt.Printf("Setting `%s/%s` to `%d`\n", environment, name, id)
	if _, ok := idMap[environment]; ok {
		idMap[environment][name] = id
	} else {
		idMap[environment] = map[string]int{name: id}
	}
	return http.StatusOK, strconv.Itoa(idMap[environment][name])
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/list", func(context *gin.Context) {
		mutex.Lock()
		context.JSON(http.StatusOK, gin.H{"idMap": idMap})
		mutex.Unlock()
	})

	router.GET("/getter/:environment/:name", func(context *gin.Context) {
		mutex.Lock()
		status, id := getID(context.Param("name"), context.Param("environment"))
		mutex.Unlock()
		context.String(status, id)
	})

	router.POST("/setter", func(context *gin.Context) {
		passedID, err := strconv.Atoi(context.PostForm("id"))
		if err != nil {
			context.String(http.StatusBadRequest, fmt.Sprintf("Error converting %s to an integer", context.PostForm("id")))
			return
		}
		mutex.Lock()
		status, id := setID(context.PostForm("name"), context.PostForm("environment"), passedID)
		mutex.Unlock()
		context.String(status, id)
	})

	return router
}

func main() {
	router := SetupRouter()
	router.Run("localhost:8080")
}
