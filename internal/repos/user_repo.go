// internal/repos/user_repo.go
package repos

import (
	"chat-back/internal/models"
	"context"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
	}
}

// ========== Métodos que ya tenía tu compañero ==========

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) SearchByUsername(ctx context.Context, query string, limit int64) ([]models.User, error) {
	pattern := "^" + regexp.QuoteMeta(query)

	filter := bson.M{
		"username": bson.M{
			"$regex":   pattern,
			"$options": "i",
		},
	}
	findOptions := options.Find().SetLimit(limit).SetSort(bson.D{{Key: "username", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) FindById(ctx context.Context, id string) (*models.User, error) {
	// Convertir el string a ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindAll(ctx context.Context) ([]models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ========== Métodos nuevos para Auth ==========

func (r *UserRepo) FindByFirebaseUID(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"firebase_uid": uid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) ExistsByFirebaseUID(ctx context.Context, uid string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"firebase_uid": uid})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update actualiza el username, birthdate o avatar_url del usuario
func (r *UserRepo) Update(ctx context.Context, id string, username *string, birthdate *time.Time, avatarURL *string) (*models.User, error) {
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
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": update}, opts).Decode(&updatedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, err
	}

	return &updatedUser, nil
}
