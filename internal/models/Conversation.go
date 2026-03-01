package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID 					primitive.ObjectID 	`bson:"_id,omitempty"`
	Type 				string             	`bson:"type"`
	Members				[]string		   	`bson:"members"`
	Name 				string             	`bson:"name,omitempty"`
	CreatedAt 			time.Time          	`bson:"created_at"`
	UpdatedAt 			time.Time          	`bson:"updated_at"`
	LastMessageAt 		*time.Time			`bson:"last_message_at,omitempty"`
	LastMessagePreview	string				`bson:"last_message_preview,omitempty"`
}