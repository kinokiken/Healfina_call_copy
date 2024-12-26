package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DialogSummary struct {
	DialogID string `bson:"dialog_id" json:"dialog_id"`
	Summary  string `bson:"summary" json:"summary"`
	Time     string `bson:"time" json:"time"`
}

type Summary struct {
	UserID          int             `bson:"_id" json:"user_id"`
	DialogSummaries []DialogSummary `bson:"dialog_summaries" json:"dialog_summaries"`
	OverallSummary  string          `bson:"overall_summary" json:"overall_summary"`
}

func AddInteractionSummary(userID int, summary string, overallSummary string) error {
	collection := Client.Database("chatgpt_telegram_bot").Collection("summary")

	timeNow := time.Now().In(time.UTC).Format("2006-01-02 15:04:05")

	currentDialogID, err := GetUserAttribute(userID, "current_dialog_id")
	if err != nil {
		return fmt.Errorf("failed to get current dialog ID: %v", err)
	}

	dialogID, ok := currentDialogID.(string)
	if !ok || dialogID == "" {
		return fmt.Errorf("current dialog ID is not valid")
	}

	dialogSummary := DialogSummary{
		DialogID: dialogID,
		Summary:  summary,
		Time:     timeNow,
	}

	filter := bson.M{"_id": userID}
	var existingSummary Summary

	err = collection.FindOne(context.Background(), filter).Decode(&existingSummary)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("error finding user summary: %v", err)
	}

	if existingSummary.UserID == 0 {
		newSummary := Summary{
			UserID:          userID,
			DialogSummaries: []DialogSummary{dialogSummary},
			OverallSummary:  overallSummary,
		}

		_, err := collection.InsertOne(context.Background(), newSummary)
		if err != nil {
			return fmt.Errorf("error inserting new summary: %v", err)
		}
	} else {
		update := bson.M{
			"$push": bson.M{"dialog_summaries": dialogSummary},
			"$set":  bson.M{"overall_summary": overallSummary},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return fmt.Errorf("error updating summary: %v", err)
		}
	}

	return nil
}

func GetOverallSummary(userID int) (string, error) {
	collection := Client.Database("chatgpt_telegram_bot").Collection("summary")

	filter := bson.M{"_id": userID}

	var summary Summary

	err := collection.FindOne(context.Background(), filter).Decode(&summary)
	fmt.Println(err)
	fmt.Println(userID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no summary found for user with ID %d", userID)
		}
		return "", fmt.Errorf("error fetching overall summary for user %d: %v", userID, err)
	}

	return summary.OverallSummary, nil
}

func GetLatestDialogSummary(userID int) (string, error) {
	collection := Client.Database("chatgpt_telegram_bot").Collection("summary")

	filter := bson.M{"_id": userID}

	var summary Summary

	err := collection.FindOne(context.Background(), filter).Decode(&summary)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", fmt.Errorf("error fetching summary for user %d: %v", userID, err)
	}

	if len(summary.DialogSummaries) == 0 {
		return "", nil
	}

	latestSummary := summary.DialogSummaries[len(summary.DialogSummaries)-1]
	return latestSummary.Summary, nil
}
