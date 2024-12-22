package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecordAfter struct {
	ID       primitive.Binary `bson:"_id"`
	UserID   int              `bson:"user_id"`
	DialogID string           `bson:"dialog_id"`
	Record   interface{}      `bson:"record"`
	Time     string           `bson:"time"`
}

func AddRecordAfter(userID int, recordAfter interface{}) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	currentDialogID, err := GetUserAttribute(userID, "current_dialog_id")
	if err != nil {
		return fmt.Errorf("failed to get current dialog ID: %v", err)
	}

	dialogIDStr, ok := currentDialogID.(string)
	if !ok || dialogIDStr == "" {
		return fmt.Errorf("current dialog ID is not valid")
	}

	recordID := uuid.New().String()
	recordDoc := RecordAfter{
		ID:       primitive.Binary{Subtype: 4, Data: []byte(recordID)},
		UserID:   userID,
		DialogID: dialogIDStr,
		Record:   recordAfter,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}

	recordCollection := Client.Database("chatgpt_telegram_bot").Collection("record")
	_, err = recordCollection.InsertOne(context.Background(), recordDoc)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}

	dialogCollection := Client.Database("chatgpt_telegram_bot").Collection("dialog")
	update := bson.M{
		"$push": bson.M{
			"records_after": fmt.Sprintf("Дата и время: %s: \n%v", recordDoc.Time, recordAfter),
		},
	}

	filter := bson.M{"_id": dialogIDStr, "user_id": userID}
	_, err = dialogCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update dialog records: %v", err)
	}

	return nil
}

func GetUserRecords(userID int) ([]RecordAfter, error) {
	recordCollection := Client.Database("chatgpt_telegram_bot").Collection("record")

	_, err := GetUserByID(int(userID))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("error checking user existence: %v", err)
	}

	filter := bson.M{"user_id": userID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := recordCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching user records: %v", err)
	}
	defer cursor.Close(ctx)

	var records []RecordAfter
	for cursor.Next(ctx) {
		var record RecordAfter
		if err := cursor.Decode(&record); err != nil {
			return nil, fmt.Errorf("error decoding record: %v", err)
		}
		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records, nil
}

func UpdateUserRecord(userID int, recordID string, updatedData interface{}) error {
	recordCollection := Client.Database("chatgpt_telegram_bot").Collection("record")

	recordUUID, err := uuid.Parse(recordID)
	if err != nil {
		return fmt.Errorf("invalid record ID format: %v", err)
	}
	recordBinary := primitive.Binary{Subtype: 4, Data: recordUUID[:]}

	filter := bson.M{"_id": recordBinary, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"record": updatedData,
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	result, err := recordCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("record not found for update")
	}

	return nil
}

func DeleteUserRecord(userID int, recordID string) error {
	recordCollection := Client.Database("chatgpt_telegram_bot").Collection("record")

	recordUUID, err := uuid.Parse(recordID)
	if err != nil {
		return fmt.Errorf("invalid record ID format: %v", err)
	}
	recordBinary := primitive.Binary{Subtype: 4, Data: recordUUID[:]}

	filter := bson.M{"_id": recordBinary, "user_id": userID}

	result, err := recordCollection.DeleteOne(context.Background(), filter)
	fmt.Println(result)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("record not found or already deleted")
	}

	return nil
}
