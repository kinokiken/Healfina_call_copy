package database

import (
	"context"
	"fmt"
	"time"

	"Healfina_call/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Dialog struct {
	ID           string    `bson:"_id" json:"id"`
	UserID       int       `bson:"user_id" json:"user_id"`
	Plan         string    `bson:"plan" json:"plan"`
	ChatMode     string    `bson:"chat_mode" json:"chat_mode"`
	StartTime    time.Time `bson:"start_time" json:"start_time"`
	Model        string    `bson:"model" json:"model"`
	Summary      string    `bson:"summary" json:"summary"`
	Messages     []string  `bson:"messages" json:"messages"`
	RecordsAfter []string  `bson:"records_after" json:"records_after"`
	UserRate     *int8     `bson:"user_rate" json:"user_rate"`
}

func AddDialog(userID int, plan string) (string, error) {
	dialogCollection := Client.Database("chatgpt_telegram_bot").Collection("dialog")
	userCollection := Client.Database("chatgpt_telegram_bot").Collection("user")

	user, err := GetUserByID(int(userID))
	if err != nil {
		return "", fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}
	current_model, err := GetUserAttribute(int(userID), "current_model")
	if err != nil {
		return "", fmt.Errorf("failed to get current model for user: %v", err)
	}

	dialogID := uuid.New().String()
	dialog := Dialog{
		ID:           dialogID,
		UserID:       userID,
		Plan:         plan,
		ChatMode:     "assistant",
		StartTime:    time.Now(),
		Model:        current_model.(string),
		Summary:      "",
		Messages:     []string{},
		RecordsAfter: []string{},
		UserRate:     nil,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = dialogCollection.InsertOne(ctx, dialog)
	if err != nil {
		return "", fmt.Errorf("failed to insert dialog: %v", err)
	}

	filter := bson.M{"chat_id": userID}
	update := bson.M{
		"$set": bson.M{
			"current_dialog_id": dialogID,
			"state":             "session",
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", fmt.Errorf("failed to update user: %v", err)
	}

	return dialogID, nil
}

func SetDialogMessages(userID int, dialogMessages []*models.SetMessagesRequestMessagesItems0, dialogID *string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if dialogID == nil || *dialogID == "" {
		currentDialogID, err := GetUserAttribute(userID, "current_dialog_id")
		if err != nil {
			return fmt.Errorf("failed to get current dialog ID: %v", err)
		}
		if dialogIDStr, ok := currentDialogID.(string); ok {
			dialogID = &dialogIDStr
		} else {
			return fmt.Errorf("current dialog ID is not a string")
		}
	}

	dialogCollection := Client.Database("chatgpt_telegram_bot").Collection("dialog")

	formattedMessages := make([]bson.M, len(dialogMessages))
	for i, msg := range dialogMessages {
		formattedMessages[i] = bson.M{
			"user": msg.User,
			"bot":  msg.Bot,
			"date": msg.Date,
		}
	}

	filter := bson.M{"_id": *dialogID, "user_id": userID}

	update := bson.M{
		"$push": bson.M{
			"messages": bson.M{
				"$each": formattedMessages,
			},
		},
	}

	_, err = dialogCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update dialog messages: %v", err)
	}

	return nil
}
