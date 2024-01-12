package mongodb

import (
	"context"
	"errors"
	"github.com/degeboman/gas/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type UsersStorage struct {
	users Users
}

func (u UsersStorage) UserByEmail(email string) (interface{}, error) {
	res := u.users.FindOne(context.TODO(), bson.D{
		{Key: "email", Value: email},
	})

	var userInfo map[string]interface{}

	if err := res.Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

type Users struct {
	*mongo.Collection
}

func (u UsersStorage) DoesEmailExist(email string) error {
	res := u.users.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil
		}

		return res.Err()
	}

	return errors.New("email is already in use")
}

func (u UsersStorage) CreateUser(email, password string, userInfo interface{}) (string, error) {
	res, err := u.users.InsertOne(context.TODO(), bson.D{
		{Key: "email", Value: email},
		{Key: "password", Value: password},
		{Key: "user_info", Value: userInfo},
	})

	if err != nil {
		return "0", err
	}

	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func New(connectString string) UsersStorage {
	const op = "storage.mongodb.mongodb.New"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectString))

	if err != nil {
		log.Fatalf("%s: %s", op, err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("%s: %s", op, err)
	}

	users := Users{
		Collection: client.Database(constant.DatabaseName).Collection("users"),
	}

	return UsersStorage{
		users: users,
	}
}
