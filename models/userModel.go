package models


type User struct {
	ID              uint8
	Name            string
	Position        string
	Age             uint
	Approval_Status string `gorm:"default:Pending"`
}
