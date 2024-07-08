package main

import (
	"github/Sadotib/go-crud/initializers"
	"github/Sadotib/go-crud/models"
)
func init(){
	initializers.ConnectDatabase()
}
func main() {
	initializers.DB.AutoMigrate(&models.User{},  &models.Events{},&models.Registration{})
}
