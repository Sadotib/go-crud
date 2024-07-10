package handlers

import (
	"fmt"
	"github/Sadotib/go-crud/globals"
	"github/Sadotib/go-crud/initializers"
	"github/Sadotib/go-crud/models"
	"os"

	"time"
	"unicode"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
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
	formData := make(map[string]string)
	//FormValue retrieves a single value by key, works for both GET and POST requests
	formData["user"] = c.Request.FormValue("username")
	formData["password"] = c.Request.FormValue("pass")
	fmt.Println(formData["user"])

	passUINT, _ := strconv.Atoi(formData["password"])

	if username != formData["user"] || pswrd != passUINT {

		fmt.Println(passUINT)

		c.Redirect(http.StatusFound, "/login?admin_privileges=yes")

		return
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": user,
	// 	"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	// })
	tokenLogin := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": formData["user"],
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	fmt.Print("bomb")
	// tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
	// if err != nil {
	// 	http.Error(c.Writer, "Error creating token", http.StatusBadRequest)
	// }
	tokenLoginString, err := tokenLogin.SignedString([]byte(os.Getenv("SECRET_TOKEN_ADMIN")))
	if err != nil {
		http.Error(c.Writer, "Error generating login token", http.StatusInternalServerError)
		return
	}

	// c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600*24*7, "", "", false, true)
	// c.SetCookie("jwt", tokenString, 3600*24, "", "", true, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "admintoken",
		Value:    tokenLoginString,
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

func AdminDashboard(c *gin.Context) {
	//adm, _ := c.Get("admin")

	tokenCookie, err := c.Request.Cookie("admintoken")

	if err != nil {
		http.Error(c.Writer, "Invalid token 1", http.StatusUnauthorized)

		return
	}
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_TOKEN_ADMIN")), nil
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
			Name:   "admintoken",
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
			err := globals.TPL.ExecuteTemplate(c.Writer, "success.html", struct {
				Message string
			}{
				Message: "Values Updated",
			})
			if err != nil {
				http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
				return
			}
		} else if action == "reject" {
			// perform reject action for row with ID row.ID
			initializers.DB.Model(&models.User{}).Where("id = ?", peep.ID).Update("approval_status", "Rejected")
			err := globals.TPL.ExecuteTemplate(c.Writer, "success.html", struct {
				Message string
			}{
				Message: "Values Updated",
			})
			if err != nil {
				http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
				return
			}
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
	err := globals.TPL.ExecuteTemplate(c.Writer, "success.html", struct {
		Message string
	}{
		Message: "Event Creation Successful!",
	})
	if err != nil {
		http.Error(c.Writer, "Error rendering template", http.StatusInternalServerError)
		return
	}
	//c.JSON(http.StatusCreated, gin.H{"message": "event created successfully"})
}

func RegisterUser(c *gin.Context) {
	var user models.User

	// Parse the form data. PostForm returns a map of form values and only works for POST requests

	name := c.PostForm("name")
	email := c.PostForm("email")
	age := c.PostForm("age")
	password := c.PostForm("password")
	guid := xid.New()

	ageUINT, _ := strconv.Atoi(age) //convert string to int using strconv
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 7 < len(password) && len(password) < 60 {
		pswdLength = true
	}

	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces {
		globals.TPL.ExecuteTemplate(c.Writer, "userRegisterForm.html", "Password doesn't obey specified criteria")
		return
	}
	var exists bool
	t := initializers.DB.Raw("SELECT 1 FROM users WHERE name=?", name).Scan(&exists)

	if t.Error != nil {
		return
	}
	if exists {
		globals.TPL.ExecuteTemplate(c.Writer, "userRegisterForm.html", "Username already exists")
		return // username already exists
	}
	var hash []byte
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		globals.TPL.ExecuteTemplate(c.Writer, "userRegisterForm.html", "There was a problem registering")
	}

	user.ID = guid.String()
	user.Name = name
	user.Email = email
	user.Password = string(hash)
	user.Age = uint8(ageUINT)
	globals.X = user.ID

	users := models.User{ID: user.ID, Name: user.Name, Email: user.Email, Age: user.Age, Password: user.Password}
	result := initializers.DB.Create(&users)

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

var u string

func LoginUser(c *gin.Context) {
	username := c.PostForm("name")
	password := c.PostForm("password")
	u = username
	var hash string
	t := initializers.DB.Raw("SELECT password FROM users WHERE name=?", username).Scan(&hash)
	if t.Error != nil {
		globals.TPL.ExecuteTemplate(c.Writer, "userLoginForm.html", "Check username and password")

		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// returns nill on succcess
	if err != nil {
		globals.TPL.ExecuteTemplate(c.Writer, "userLoginForm.html", "Check username and password")
		return
	}
	fmt.Println(username)
	tokenLogin := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenLoginString, err := tokenLogin.SignedString([]byte(os.Getenv("SECRET_TOKEN_USER")))
	if err != nil {
		http.Error(c.Writer, "Error generating login token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "usertoken",
		Value:    tokenLoginString,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HttpOnly: true,
	})

	c.HTML(http.StatusOK, "userDash.html", nil)

}

type Profile struct {
	ID         string
	Username   string
	Email      string
	Events     []uint8
	EventNames []string
}

func UserDashboard(c *gin.Context) {
	tokenCookie, err := c.Request.Cookie("usertoken")

	if err != nil {
		http.Error(c.Writer, "Invalid token 1", http.StatusUnauthorized)
		fmt.Println(err)

		return
	}
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_TOKEN_USER")), nil
	})
	if err != nil {
		http.Error(c.Writer, "Invalid token 2", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user"] != u {
		http.Error(c.Writer, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	action := c.Query("action")

	var b Profile
	if action == "profile" {

		initializers.DB.Raw("SELECT id FROM users WHERE name=?", u).Scan(&b.ID)
		initializers.DB.Raw("SELECT email FROM users WHERE name=?", u).Scan(&b.Email)
		initializers.DB.Raw("SELECT event_id FROM registrations WHERE user_id=?", b.ID).Scan(&b.Events)
		// for i := 0; i < len(b.Events); i++ {
		// 	k := b.Events[i]
		// 	var p string
		// 	initializers.DB.Raw("SELECT event_name FROM events WHERE event_id=?",k).Scan(&p)
		// 	b.EventNames=append(b.EventNames, p)
		// }
		for _, k := range b.Events {
			var p string
			initializers.DB.Raw("SELECT event_name FROM events WHERE event_id=?", k).Scan(&p)
			b.EventNames = append(b.EventNames, p)
		}
		b.Username = u
		fmt.Println(b.EventNames)
		fmt.Println("kkkk")

		globals.TPL.ExecuteTemplate(c.Writer, "fetchProfile.html", b)

	} else if action == "logout" {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:   "usertoken",
			Value:  "",
			MaxAge: -1,
		})
		c.Redirect(http.StatusFound, "/action?action=login")
	}

}
