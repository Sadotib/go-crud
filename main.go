package main

import (
	"github/Sadotib/go-crud/globals"
	"github/Sadotib/go-crud/handlers"
	"github/Sadotib/go-crud/initializers"
	"log"
	"os"

	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.ConnectDatabase()
}

func main() {
	r := gin.Default()

	globals.TPL, _ = template.ParseGlob("templates/*.html")
	// r.LoadHTMLFiles("index.html")
	// r.LoadHTMLGlob("templates/*.html")
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	} 

	t, err = t.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminYN.html", nil)
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "templates/favicon.ico")
	})

	r.GET("/login", handlers.AdminPriv)

	r.POST("/admin", handlers.AdminLogin)

	r.GET("/admindash", handlers.AdminDashboard)

	r.POST("/admindash/entries", handlers.AcceptReject)

	r.POST("/admindash/create", handlers.CreateEvent)

	r.GET("/action", func(c *gin.Context) {
		// TODO: will do later
		create := c.Query("action")

		if create == "register" {
			c.HTML(http.StatusOK, "userRegisterForm.html", nil)
		} else if create == "login" {
			//write code here
			c.HTML(http.StatusOK, "userLoginForm.html", nil)
		}
	})

	r.POST("/user", handlers.LoginUser)

	r.POST("/userregister", handlers.RegisterUser)

	r.POST("/userregister/selectevent", handlers.SelectNRegister)
	r.GET("/userdash", handlers.UserDashboard)

	port := os.Getenv("PORT")
	if port ==""{
		port= "3000"
	}

	r.Run("0.0.0.0:" + port)
}

//TODO: 1) fix the foreign key problem and registration db not getting created [FIXED]
//2) fix the return message after succesfull user or event creation
//3) watch the video on keeping users logged in which will be used to keep admin logged in
//4) if possible, return the data entered by user after successful registration with their info and events they have registered for
//5) make a table for the users using which they can check all the events available and register for each one using buttons [FIXED]
//6) make a home page before asking if you are admin. also give functionalities to the navigation menu buttons
//7) fix the fetchEvents table, update the columns [FIXED]
//8) if possible, give an email feature using that tool where user recieves email notification of successful registration
//9) do some error handling and give a creative "page not found" page
//10) fix the approve and reject of the fetchEntries file so that admin can approve or reject. If rejected, delete the user and if possible notify them using email
