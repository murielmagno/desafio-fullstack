package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Person struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  string    `json:"birthday"`
	Password  string    `json:"password"`
	UserName  string    `json:"username" gorm:"unique"`
	Friends   []*Person `gorm:"many2many:person_friends;association_jointable_foreignkey:friend_id"` // Lista de amigos
	Cards     []*Card   `gorm:"foreignKey:user_id"`
}

type Card struct {
	ID           uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"card_id"`
	Title        string    `json:"title"`
	PAN          string    `json:"pan" gorm:"unique"`
	ExpiryMM     string    `json:"expiry_mm"`
	ExpiryYYYY   string    `json:"expiry_yyyy"`
	SecurityCode string    `json:"security_code"`
	Limit        float64   `json:"limit" gorm:"not null;default:0"`
	Balance      float64   `json:"balance" gorm:"not null;default:0"`
	Date         time.Time `json:"date"`
	UserId       uuid.UUID
}

type Transfer struct {
	SenderID        string  `json:"sender_id"`
	FriendID        string  `json:"friend_id"`
	TotalToTransfer float64 `json:"total_to_transfer"`
	Pan             string
	CardId          string `json:"card_id"`
}
