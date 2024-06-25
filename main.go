package main

import (
	"fmt"
	"github/Sadotib/go-crud/initializers"
	"github/Sadotib/go-crud/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	 
)
func init() {
	initializers.ConnectDatabase()
}
func main() {

	initializers.DB.AutoMigrate(&models.User{})

	r := gin.Default()
	r.LoadHTMLGlob("forms/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "adminYN.html", nil)

	})

	r.GET("/admin", func(c *gin.Context) {
		admin_priv := c.Query("admin_privileges")
		// c.String(http.StatusOK, "Hello %s", admin_priv)
		if admin_priv == "yes" {
			c.HTML(http.StatusOK, "adminForm.html", nil)
		}

		if admin_priv == "no" {

			c.HTML(http.StatusOK, "mainForm.html", nil)

		}
	})

	r.POST("/admin", func(c *gin.Context) {
		username := "admin"
		pass := 12345

		c.Request.ParseForm()

		//FormValue retrieves a single value by key, works for both GET and POST requests
		user := c.Request.FormValue("username")
		password := c.Request.FormValue("pass")

		passUINT, _ := strconv.Atoi(password)
		var people []models.User
		if username == user && pass == passUINT {

			c.HTML(http.StatusOK, "adminApr.html", nil)

			initializers.DB.Find(&people)

			for _, i := range people {
				fmt.Printf("|ID: %d | Name: %s | Position: %s | Age: %d | Approval Status: %s |\n", i.ID, i.Name, i.Position, i.Age, i.Approval_Status)
			}

			fmt.Print("\nApprove RegistrationID (0 for NONE): ")
			var ag int
			fmt.Scan(&ag)

			if ag == 0 {
				fmt.Print("\nNO APPROVAL GRANTED\n")
			} else {
				initializers.DB.Model(&people).Where("ID = ?", ag).Update("Approval_Status", "Granted")
			}
			fmt.Print("\nDeny RegistrationID (0 for NONE): ")
			var ad int
			fmt.Scan(&ad)

			if ad == 0 {
				fmt.Print("\nNO APPROVAL DENIED\n\n")
			} else {
				initializers.DB.Model(&people).Where("ID = ?", ag).Update("Approval_Status", "Denied")
			}

		}
		if username != user || pass != passUINT {

			c.String(http.StatusOK, "Incorrect Username or Password")
		}
	})
	r.GET("/submit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user_form.html", nil)
	})
	r.POST("/submit", func(c *gin.Context) {
		var employee models.User

		// Parse the form data. PostForm returns a map of form values and only works for POST requests
		id := c.PostForm("id")
		name := c.PostForm("name")
		position := c.PostForm("position")
		age := c.PostForm("age")

		ageUINT, _ := strconv.Atoi(age) //convert string to int using strconv
		idUINT, _ := strconv.Atoi(id)   //convert string to int using strconv

		fmt.Println(idUINT)
		employee.ID = uint8(idUINT)
		fmt.Println(employee.ID)
		employee.Name = name
		employee.Position = position
		employee.Age = uint(ageUINT)

		// err := c.Bind(&employee)
		// if err!= nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		user := models.User{ID: employee.ID, Name: employee.Name, Position: employee.Position, Age: employee.Age}
		initializers.DB.Create(&user)

		c.JSON(http.StatusCreated, gin.H{"message": "Employee created successfully"})
	})

	r.Run(":8080")
}
