package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID  `json:"id"`
	UserId    uuid.UUID  `json:"user_id"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type User struct {
	Id   uuid.UUID
	Name string
}
