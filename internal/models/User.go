package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	FirebaseUID string             `bson:"firebase_uid"   json:"-"` // ← NUEVO (json:"-" para no exponerlo)
	Email       string             `bson:"email"          json:"email"`
	Username    string             `bson:"username"       json:"username"`
	AvatarURL   string             `bson:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"     json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"     json:"updated_at"`
}
