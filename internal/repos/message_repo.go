package repos

import (
	"chat-back/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepo struct {
	collection *mongo.Collection
}

func NewMessageRepo(db *mongo.Database) *MessageRepo {
	return &MessageRepo{collection: db.Collection("messages")}
}

func (r *MessageRepo) Create(ctx context.Context, msg *models.Message) error {
	msg.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, msg)
	return err
}

// Obtener mensajes con paginación por cursor (más eficiente que skip/limit)
func (r *MessageRepo) FindByConversation(
	ctx context.Context,
	convID primitive.ObjectID,
	limit int64,
	before *primitive.ObjectID, // cursor para paginación
) ([]models.Message, error) {
	filter := bson.M{"conversation_id": convID, "deleted": false}
	if before != nil {
		filter["_id"] = bson.M{"$lt": *before}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	return messages, cursor.All(ctx, &messages)
}
