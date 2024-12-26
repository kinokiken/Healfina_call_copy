package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserID              int        `bson:"chat_id" json:"user_id"`
	FirstName           string     `bson:"first_name" json:"first_name"`
	LastName            *string    `bson:"last_name" json:"last_name"`
	NSessions           int8       `bson:"n_sessions" json:"n_sessions"`
	NUnrecordedSessions int8       `bson:"n_unrecorded_sessions" json:"n_unrecorded_sessions"`
	SubscriptionStart   *time.Time `bson:"subscription_start" json:"subscription_start"`
	SubscriptionEnd     *time.Time `bson:"subscription_end" json:"subscription_end"`
	ChatMode            string     `bson:"chat_mode" json:"chat_mode"`
	Password            *string    `bson:"password,omitempty"`
}

func SetUserPassword(userID int, password string) error {
	collection := Client.Database("chatgpt_telegram_bot").Collection("user")

	filter := bson.M{"chat_id": userID}

	update := bson.M{
		"$set": bson.M{"password": password}, // в реальном проекте -> хеш
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found or already has password")
	}
	return nil
}

func GetUserByID(user_id int) (*User, error) {

	collection := Client.Database("chatgpt_telegram_bot").Collection("user")

	var user User

	filter := bson.M{"chat_id": user_id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := collection.FindOne(ctx, filter).Decode(&user)
	if data != nil {
		fmt.Println(data)
		if data == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("Ошибка при получении пользователя: %v", data)
		return nil, data
	}

	return &user, nil
}

func GetUserAttribute(user_id int, attribute string) (interface{}, error) {
	collection := Client.Database("chatgpt_telegram_bot").Collection("user")

	filter := bson.M{"chat_id": user_id}

	var result bson.M

	projection := bson.M{attribute: 1, "_id": 0}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with ID %d not found", user_id)
		}
		return nil, fmt.Errorf("error fetching attribute: %v", err)
	}

	value, exists := result[attribute]
	if !exists {
		return nil, fmt.Errorf("attribute %s not found for user %d", attribute, user_id)
	}

	return value, nil
}

func SetUserAttribute(user_id int, attribute string, value interface{}) error {
	collection := Client.Database("chatgpt_telegram_bot").Collection("user")

	filter := bson.M{"chat_id": user_id}

	update := bson.M{"$set": bson.M{attribute: value}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user with ID %d not found", user_id)
		}
		return fmt.Errorf("error updating attribute: %v", err)
	}

	log.Printf("Attribute %s successfully updated for user %d", attribute, user_id)
	return nil
}
