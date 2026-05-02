package repos

import (
	"chat-back/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Update actualiza el username, birthdate o avatar_url del usuario
func UpdateUser(collection *mongo.Collection) func(ctx context.Context, id string, username *string, birthdate *time.Time, avatarURL *string) (*models.User, error) {
	return func(ctx context.Context, id string, username *string, birthdate *time.Time, avatarURL *string) (*models.User, error) {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("ID inválido: %w", err)
		}

		update := bson.M{
			"updated_at": time.Now(),
		}

		// Solo actualizar los campos que se proporcionaron
		if username != nil && *username != "" {
			update["username"] = *username
		}
		if birthdate != nil {
			update["birthdate"] = birthdate
		}
		if avatarURL != nil && *avatarURL != "" {
			update["avatar_url"] = *avatarURL
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updatedUser models.User
		err = collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": update}, opts).Decode(&updatedUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("usuario no encontrado")
			}
			return nil, err
		}

		return &updatedUser, nil
	}
}
