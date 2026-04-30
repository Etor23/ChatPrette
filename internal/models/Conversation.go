package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Type          string               `bson:"type" json:"type"` // "dm" o "group"
	Name          string               `bson:"name,omitempty" json:"name,omitempty"`
	Members       []primitive.ObjectID `bson:"members" json:"members"`
	CreatedBy     primitive.ObjectID   `bson:"created_by" json:"created_by"`
	LastMessageAt time.Time            `bson:"last_message_at" json:"last_message_at"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
}
