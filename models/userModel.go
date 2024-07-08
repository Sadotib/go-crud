package models

type User struct {
	ID              string         `gorm:"primaryKey;not null;autoIncrement:false"`
	Name            string         `gorm:"not null"`
	Position        string         `gorm:"not null"`
	Age             uint8          `gorm:"not null"`
	Approval_Status string         `gorm:"default:Pending"`
	Registrations   []Registration `gorm:"foreignkey:UserID"`
	// TODO: fix the problem of registration table not getting created and check the foreign key problem
}

type Events struct {
	EventID       uint8          `gorm:"primaryKey;not null;autoIncrement:true"`
	EventName     string         `gorm:"unique;not null"`
	Description   string         `gorm:"not null"`
	Registrations []Registration `gorm:"foreignkey:EventID"`
}

type Registration struct {
	UserID string `gorm:"size:20"`
	// User    User `gorm:"foreignkey:UserID;references:id"`
	EventID uint8
	// 	Event   Events `gorm:"foreignkey:EventID;references:id"`
}
