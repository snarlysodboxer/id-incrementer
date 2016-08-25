package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

// TODO setup auto API documentation, add auth, add persistent storage, add settings flags or file

var initialValue = 42
var incrementBy = 5
var mutex = &sync.Mutex{}

type idMap map[string]map[string]int

func NewIDMap() idMap {
	return map[string]map[string]int{}
}

func (ids idMap) Get(name, environment string) (int, int) {
	// check if the environment is found
	if _, ok := ids[environment]; ok {
		// check if the name is found
		if id, ok := ids[environment][name]; ok {
			ids[environment][name] = id + incrementBy
		} else {
			// add unfound name
			ids[environment][name] = initialValue
		}
	} else {
		// add unfound environment and name
		ids[environment] = map[string]int{name: initialValue}
	}
	return http.StatusOK, ids[environment][name]
}

func (ids idMap) Set(name, environment string, id int) (int, int) {
	if _, ok := ids[environment]; ok {
		ids[environment][name] = id
	} else {
		ids[environment] = map[string]int{name: id}
	}
	return http.StatusOK, ids[environment][name]
}

func (ids idMap) SetupRouter() *gin.Engine {
	// log to stdout
	router := gin.Default()

	// // don't log to stdout (helpful for testing)
	// router := gin.New()
	// router.Use(gin.Recovery())

	router.GET("/lister", func(context *gin.Context) {
		mutex.Lock()
		context.JSON(http.StatusOK, ids)
		mutex.Unlock()
	})

	router.GET("/getter/:environment/:name", func(context *gin.Context) {
		mutex.Lock()
		status, id := ids.Get(context.Param("name"), context.Param("environment"))
		context.JSON(status, map[string]int{"id": id})
		mutex.Unlock()
	})

	router.POST("/setter", func(context *gin.Context) {
		mutex.Lock()
		if context.PostForm("id") == "" {
			context.JSON(http.StatusBadRequest, `{"error": "ID field was not passed or is empty"}`)
			mutex.Unlock()
			return
		}
		passedID, err := strconv.Atoi(context.PostForm("id"))
		if err != nil {
			message := fmt.Sprintf("Error converting `%s` to an integer", context.PostForm("id"))
			msg := map[string]string{"error": message}
			context.JSON(http.StatusBadRequest, msg)
			mutex.Unlock()
			return
		}
		status, id := ids.Set(context.PostForm("name"), context.PostForm("environment"), passedID)
		idM := map[string]int{"id": id}
		context.JSON(status, idM)
		mutex.Unlock()
	})

	return router
}

func main() {
	ids := NewIDMap()
	router := ids.SetupRouter()
	router.Run("localhost:8080")
}
