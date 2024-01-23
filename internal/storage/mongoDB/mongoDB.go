package mongoDB

import (
	"context"

	"github.com/IskanderSh/taqwa-auth/internal/config"
	"github.com/IskanderSh/taqwa-auth/internal/domain/models"
	"github.com/IskanderSh/taqwa-auth/internal/lib/error/wrapper"
	"github.com/IskanderSh/taqwa-auth/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client *mongo.Client
}

// New creates a new instance of the mongoDB storage
func New(db *config.DB) (*Storage, error) {
	const op = "storage.mongodb.New"
	var ctx = context.TODO()

	// TODO: replace localhost
	uri := "mongodb://localhost:27017/"
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return &Storage{}, wrapper.Wrap(op, err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return &Storage{}, wrapper.Wrap(op, err)
	}

	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, hashPass []byte) (string, error) {
	const op = "storage.SaveUser"

	coll := s.client.Database("db").Collection("users")

	user := storage.User{Email: email, HashPass: hashPass}

	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	userID := result.InsertedID.(primitive.ObjectID).String()

	return userID, nil
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.User"

	coll := s.client.Database("db").Collection("users")

	filter := bson.D{{"email", email}}

	var user models.User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, storage.ErrUserNotFound
	}

	return &user, nil
	//
	//return mapToUserModel(user), nil
}

//func mapToUserModel(user storage.User) *models.User {
//	return &models.User{
//		Email:    user.Email,
//		HashPass: user.HashPass,
//	}
//}
