package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
)

func main() {
	// Default returns an Engine instance with the Logger and some middleware already attached. We can also use the
	//New method which is the same but the middleware is absent
	// func Default() *Engine
	router := gin.Default()

	router.LoadHTMLGlob("hello/*.html") //we can also use LoadHTMLGlob such everything inside the folder that matches
	//the pattern will load. Here we are using just one html file

	router.GET("/hello", getHello)

	router.GET("/greet", getGreet)

	// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
	// func (engine *Engine) Run(addr ...string) (err error)

	router.Run("localhost:9999")

}

// GET /hello
func getHello(c *gin.Context) {
	// String writes the given string into the response body.
	// http.StautsOK is http status code saved as constant in http package
	// func (c *Context) String(code int, format string, values ...interface{})
	c.String(http.StatusOK, "ok",nil)
}

// GET /greet
func getGreet(c *gin.Context) {
	// HTML renders the HTTP template specified by its file name.
	// It also updates the HTTP code and sets the Content-Type as "text/html"
	// func (c *Context) HTML(code int, name string, obj interface{})
	c.HTML(200, "hello.html", nil)
}

