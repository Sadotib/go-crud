package main

import (
	"github/Sadotib/go-crud/globals"
	"github/Sadotib/go-crud/handlers"
	"github/Sadotib/go-crud/initializers"
	

	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.ConnectDatabase()
}

func main() {
	globals.TPL, _ = template.ParseGlob("templates/*.html")
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminYN.html", nil)
	})

	r.GET("/login", handlers.AdminPriv)

	r.POST("/admin", handlers.AdminLogin)

	// r.GET("/validate", middleware.RequireAuth,handlers.Validate)

	r.GET("/admin", handlers.AdminDashboard)

	r.POST("/admin/entries", handlers.AcceptReject)

	r.POST("/admin/create", handlers.CreateEvent)

	r.GET("/create", func(c *gin.Context) {
		// TODO: will do later
		create := c.Query("create")

		if create == "user" {
			c.HTML(http.StatusOK, "userForm.html", nil)
		}
	})

	r.POST("/create/user", handlers.EnterUserDetails)

	r.POST("/create/user/selectevent", handlers.SelectNRegister)

	r.Run(":8080")
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
