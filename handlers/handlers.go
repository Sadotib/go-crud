package handlers

import (
	"fmt"
	"github/Sadotib/go-crud/globals"
	"github/Sadotib/go-crud/initializers"
	"github/Sadotib/go-crud/models"
	"os"
	"time"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
)

func AdminPriv(c *gin.Context) {
	admin_priv := c.Query("admin_privileges")
	// c.String(http.StatusOK, "Hello %s", admin_priv)
	if admin_priv == "yes" {
		c.HTML(http.StatusOK, "adminForm.html", nil)
	}
	if admin_priv == "no" {

		c.HTML(http.StatusOK, "create.html", nil)
	}
}

func AdminLogin(c *gin.Context) {
	godotenv.Load()
	pass := os.Getenv("ADMIN_PASSWORD")
	username := os.Getenv("ADMIN_USERNAME")

	pswrd, _ := strconv.Atoi(pass)

	c.Request.ParseForm()

	//FormValue retrieves a single value by key, works for both GET and POST requests
	user := c.Request.FormValue("username")
	password := c.Request.FormValue("pass")

	passUINT, _ := strconv.Atoi(password)

	if username != user || pswrd != passUINT {

		c.Redirect(http.StatusFound, "/login?admin_privileges=yes")

		return
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": user,
	// 	"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	// })
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	fmt.Print("bomb")
	// tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
	// if err != nil {
	// 	http.Error(c.Writer, "Error creating token", http.StatusBadRequest)
	// }
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
	if err != nil {
		http.Error(c.Writer, "Error generating token", http.StatusInternalServerError)
		return
	}
	// c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)
	// c.SetCookie("jwt", tokenString, 3600*24, "", "", true, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HttpOnly: true,
	})
	c.HTML(http.StatusOK, "adminDash.html", nil)
	// http.Redirect(c.Writer, c.Request, "/admin/dashboard", http.StatusFound)
	// c.Redirect(http.StatusFound,"/admin/dashboard")
	// c.Abort()
	// c.Request.Body.Close()
	//c.Redirect(http.StatusFound, "/admin")

	// c.HTML(http.StatusOK, "adminDash.html", nil)

}

// func AdminDashboardHandler(c *gin.Context) {
//     // Verify token and render dashboard
//     token := c.GetCookie("auth_token")
//     if token != "" {
//         // Render admin dashboard HTML page
//         c.HTML(http.StatusOK, "adminDash.html", nil)
//     } else {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
//     }
// }

func AdminDashboard(c *gin.Context) {
	//adm, _ := c.Get("admin")
	tokenCookie, err := c.Request.Cookie("jwt")
	if err != nil {
		http.Error(c.Writer, "Invalid token 1", http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_TOKEN")), nil
	})
	if err != nil {
		http.Error(c.Writer, "Invalid token 2", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user"] != os.Getenv("ADMIN_USERNAME") {
		http.Error(c.Writer, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	fmt.Println("pp")

	view := c.Query("view")
	// c.String(http.StatusOK, "Hello %s", admin_priv)
	if view == "entries" {

		initializers.DB.Raw("SELECT * FROM users").Scan(&globals.PEOPLE)
		globals.TPL.ExecuteTemplate(c.Writer, "fetchEntries.html", globals.PEOPLE)

	} else if view == "events" {

		initializers.DB.Raw("SELECT * FROM events").Scan(&globals.EVENT)
		globals.TPL.ExecuteTemplate(c.Writer, "fetchEvents.html", globals.EVENT)

	} else if view == "createevent" {
		c.HTML(http.StatusOK, "eventForm.html", nil)

	} else if view == "logout" { //TODO: added this here
		http.SetCookie(c.Writer, &http.Cookie{
			Name:   "jwt",
			Value:  "",
			MaxAge: -1,
		})

		fmt.Print("hello")
		c.Redirect(http.StatusFound, "/login?admin_privileges=yes")

	}

}

func AcceptReject(c *gin.Context) {
	if c.Request.Method != "POST" {
		log.Println("Invalid request method:", c.Request.Method)
		return
	}
	for _, peep := range globals.PEOPLE { // assume rows is the slice of values
		action := c.Request.FormValue(fmt.Sprintf("action_%s", peep.ID))
		if action == "accept" {
			// perform accept action for row with ID row.ID
			initializers.DB.Model(&models.User{}).Where("id = ?", peep.ID).Update("approval_status", "Accepted")
		} else if action == "reject" {
			// perform reject action for row with ID row.ID
			initializers.DB.Model(&models.User{}).Where("id = ?", peep.ID).Update("approval_status", "Rejected")
		}
	}

}

func CreateEvent(c *gin.Context) {
	var event models.Events
	// Parse the form data. PostForm returns a map of form values and only works for POST requests
	eventname := c.PostForm("eventname")
	desc := c.PostForm("description")
	event.EventName = eventname
	event.Description = desc
	evnt := models.Events{EventName: event.EventName, Description: event.Description}
	result := initializers.DB.Create(&evnt)
	if result.Error != nil {
		http.Error(c.Writer, "Error creating event", http.StatusBadRequest)
	}
	c.JSON(http.StatusCreated, gin.H{"message": "event created successfully"})
}

func EnterUserDetails(c *gin.Context) {
	var employee models.User

	// Parse the form data. PostForm returns a map of form values and only works for POST requests

	name := c.PostForm("name")
	position := c.PostForm("position")
	age := c.PostForm("age")
	guid := xid.New()

	ageUINT, _ := strconv.Atoi(age) //convert string to int using strconv

	employee.ID = guid.String()
	employee.Name = name
	employee.Position = position
	employee.Age = uint8(ageUINT)
	globals.X = employee.ID

	user := models.User{ID: employee.ID, Name: employee.Name, Position: employee.Position, Age: employee.Age}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		http.Error(c.Writer, "Error creating user", http.StatusBadRequest)
	}
	var event []models.Events
	initializers.DB.Raw("SELECT * FROM events").Scan(&event)
	globals.TPL.ExecuteTemplate(c.Writer, "selectEvent.html", event)
}

func SelectNRegister(c *gin.Context) {
	if c.Request.Method != "POST" {
		log.Println("Invalid request method:", c.Request.Method)
		return
	}
	c.Request.ParseForm()
	eventIDs := c.Request.Form["event_ids[]"]
	log.Println("eventIDs:", eventIDs)

	// Create a new registration entry for each selected event ID
	for i := 0; i < len(eventIDs); i++ {
		eventID := eventIDs[i]

		fmt.Print(eventID)
		ui64, err := strconv.ParseUint(eventID, 10, 64)
		if err != nil {
			panic(err)
		}
		evntID := uint8(ui64)

		//check if already exists
		var exists bool
		//er := initializers.DB.QueryRow("SELECT 1 FROM registration WHERE userID = $1 AND EventID = $2", userID, eventID).Scan(&exists)
		t := initializers.DB.Raw("SELECT 1 FROM registrations WHERE user_id=? AND event_id=?", globals.X, evntID).Scan(&exists)

		if t.Error != nil {
			return
		}
		if exists {

			return // registration already exists, do nothing
		}

		evnt := models.Registration{UserID: globals.X, EventID: evntID}
		result := initializers.DB.Create(&evnt)
		if result.Error != nil {
			http.Error(c.Writer, "Error registering", http.StatusBadRequest)
		}
	}

	err := globals.TPL.ExecuteTemplate(c.Writer, "success.html", struct {
		Message string
	}{
		Message: "Registration successful!",
	})
	if err != nil {
		http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
