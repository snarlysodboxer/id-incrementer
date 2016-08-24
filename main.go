package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

// TODO setup auto API documentation, add auth, add persistent storage, add settings flags or file

type idMap map[string]map[string]int

var initialValue = 42
var incrementBy = 5
var mutex = &sync.Mutex{}

func (ids idMap) Get(name, environment string) (int, string) {
	// check if the environment is found
	if _, ok := ids[environment]; ok {
		// check if the name is found
		if id, ok := ids[environment][name]; ok {
			ids[environment][name] = id + incrementBy
		} else {
			// add unfound name
			// fmt.Printf("Adding `%s/%s` with initial value `%d`\n", environment, name, initialValue)
			ids[environment][name] = initialValue
		}
	} else {
		// add unfound environment and name
		// fmt.Printf("Adding `%s/%s` with initial value `%d`\n", environment, name, initialValue)
		ids[environment] = map[string]int{name: initialValue}
	}
	return http.StatusOK, strconv.Itoa(ids[environment][name])
}

func (ids idMap) Set(name, environment string, id int) (int, string) {
	// fmt.Printf("Setting `%s/%s` to `%d`\n", environment, name, id)
	if _, ok := ids[environment]; ok {
		ids[environment][name] = id
	} else {
		ids[environment] = map[string]int{name: id}
	}
	return http.StatusOK, strconv.Itoa(ids[environment][name])
}

func (ids idMap) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/lister", func(context *gin.Context) {
		mutex.Lock()
		context.JSON(http.StatusOK, ids)
		mutex.Unlock()
	})

	router.GET("/getter/:environment/:name", func(context *gin.Context) {
		mutex.Lock()
		status, id := ids.Get(context.Param("name"), context.Param("environment"))
		mutex.Unlock()
		context.String(status, id)
	})

	router.POST("/setter", func(context *gin.Context) {
		if context.PostForm("id") == "" {
			context.String(http.StatusBadRequest, fmt.Sprintf("`id` field was not passed or is empty"))
			return
		}
		passedID, err := strconv.Atoi(context.PostForm("id"))
		if err != nil {
			context.String(http.StatusBadRequest, fmt.Sprintf("Error converting %s to an integer", context.PostForm("id")))
			return
		}
		mutex.Lock()
		status, id := ids.Set(context.PostForm("name"), context.PostForm("environment"), passedID)
		mutex.Unlock()
		context.String(status, id)
	})

	return router
}

func NewIDMap() idMap {
	return map[string]map[string]int{}
}

func main() {
	ids := NewIDMap()
	router := ids.SetupRouter()
	router.Run("localhost:8080")
}
