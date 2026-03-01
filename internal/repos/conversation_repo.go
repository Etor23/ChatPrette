package repos

import (
	"chat-back/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationRepo struct {
	collection *mongo.Collection
}

func NewConversationRepo(db *mongo.Database) *ConversationRepo {
	return &ConversationRepo{
		collection: db.Collection("conversations"),
	}
}

func (r *ConversationRepo) Create(ctx context.Context, conversation *models.Conversation) error{

	_, err := r.collection.InsertOne(ctx, conversation)
	return err
}

func (r *ConversationRepo) FindById(ctx context.Context, id string) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepo) FindByMember(ctx context.Context, memberID string) ([]models.Conversation, error) {
	filter := bson.M{"members": memberID}
	opts := options.Find().SetSort(bson.D{{Key: "lastMessageAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}

	return conversations, nil
}

func (r *ConversationRepo) UpdateName(ctx context.Context, id string, name string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"updated_at": time.Now(),
		},
	}

	

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *ConversationRepo) AddMembers(ctx context.Context, id string, members []string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$addToSet": bson.M{
			"members": bson.M{
				"$each": members,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *ConversationRepo) RemoveMembers(ctx context.Context, id string, members []string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"$in": members,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *ConversationRepo) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}