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

type ConversationRepo struct {
	collection *mongo.Collection
}

func NewConversationRepo(db *mongo.Database) *ConversationRepo {
	return &ConversationRepo{collection: db.Collection("conversations")}
}

func (r *ConversationRepo) Create(ctx context.Context, conv *models.Conversation) error {
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()
	conv.LastMessageAt = time.Now()
	_, err := r.collection.InsertOne(ctx, conv)
	return err
}

// Buscar conversaciones donde el usuario es miembro, ordenadas por último mensaje
func (r *ConversationRepo) FindByMember(ctx context.Context, userID primitive.ObjectID) ([]models.Conversation, error) {
	opts := options.Find().SetSort(bson.D{{Key: "last_message_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"members": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var convs []models.Conversation
	return convs, cursor.All(ctx, &convs)
}

// Buscar DM existente entre dos usuarios
func (r *ConversationRepo) FindDM(ctx context.Context, userA, userB primitive.ObjectID) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.collection.FindOne(ctx, bson.M{
		"type":    "dm",
		"members": bson.M{"$all": []primitive.ObjectID{userA, userB}, "$size": 2},
	}).Decode(&conv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &conv, nil
}
