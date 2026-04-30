package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`
	SenderID       primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	Content        string             `bson:"content" json:"content"`
	Type           string             `bson:"type" json:"type"` // "text", "image", "file"
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	EditedAt       *time.Time         `bson:"edited_at,omitempty" json:"edited_at,omitempty"`
	Deleted        bool               `bson:"deleted" json:"deleted"`
}
