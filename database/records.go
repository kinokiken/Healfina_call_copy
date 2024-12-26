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

type RecordAfterResponse struct {
	ID       string      `json:"id"`
	UserID   int         `json:"user_id"`
	DialogID string      `json:"dialog_id"`
	Record   interface{} `json:"record"`
	Time     string      `json:"time"`
}

func ConvertToRecordAfterResponse(record RecordAfter) (RecordAfterResponse, error) {
	var response RecordAfterResponse
	response.UserID = record.UserID
	response.DialogID = record.DialogID
	response.Record = record.Record
	response.Time = record.Time

	rawBytes := record.ID.Data

	asciiID := string(rawBytes)
	if len(asciiID) == 36 && looksLikeUUID(asciiID) {
		// Если это UUID-строка
		if parsed, err := uuid.Parse(asciiID); err == nil {
			response.ID = parsed.String()
			return response, nil
		}
	}

	if record.ID.Subtype == 4 && len(rawBytes) == 16 {
		if decoded, err := uuid.FromBytes(rawBytes); err == nil {
			response.ID = decoded.String()
			return response, nil
		}
	}

	response.ID = asciiID
	return response, nil
}

func looksLikeUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	return s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

func AddRecordAfter(userID int, recordAfter interface{}) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to check user: %v", err)
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

	newUUID := uuid.New().String()

	recordID := primitive.Binary{
		Subtype: 4,
		Data:    []byte(newUUID),
	}

	recordDoc := RecordAfter{
		ID:       recordID,
		UserID:   userID,
		DialogID: dialogIDStr,
		Record:   recordAfter,
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}

	recordColl := Client.Database("chatgpt_telegram_bot").Collection("record")
	_, err = recordColl.InsertOne(context.Background(), recordDoc)
	if err != nil {
		return fmt.Errorf("failed to insert record: %v", err)
	}

	dialogColl := Client.Database("chatgpt_telegram_bot").Collection("dialog")
	update := bson.M{
		"$push": bson.M{
			"records_after": fmt.Sprintf("Дата и время: %s: \n%v", recordDoc.Time, recordAfter),
		},
	}

	filter := bson.M{"_id": dialogIDStr, "user_id": userID}
	_, err = dialogColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update dialog records: %v", err)
	}

	return nil
}

func GetUserRecords(userID int) ([]RecordAfterResponse, error) {
	coll := Client.Database("chatgpt_telegram_bot").Collection("record")

	_, err := GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user %d not found", userID)
		}
		return nil, err
	}

	filter := bson.M{"user_id": userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching user records: %v", err)
	}
	defer cursor.Close(ctx)

	var records []RecordAfter
	for cursor.Next(ctx) {
		var rec RecordAfter
		if err := cursor.Decode(&rec); err != nil {
			return nil, fmt.Errorf("error decoding record: %v", err)
		}
		records = append(records, rec)
	}

	var response []RecordAfterResponse
	for _, r := range records {
		rr, err := ConvertToRecordAfterResponse(r)
		if err != nil {
			return nil, err
		}
		response = append(response, rr)
	}

	if len(response) == 0 {
		return nil, nil
	}
	return response, nil
}

func UpdateUserRecord(userID int, recordID string, updatedData interface{}) error {
	coll := Client.Database("chatgpt_telegram_bot").Collection("record")

	binID := primitive.Binary{Subtype: 4, Data: []byte(recordID)}
	filter := bson.M{"_id": binID, "user_id": userID}

	update := bson.M{
		"$set": bson.M{
			"record": updatedData,
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("record not found for update")
	}
	return nil
}

func DeleteUserRecord(userID int, recordID string) error {
	coll := Client.Database("chatgpt_telegram_bot").Collection("record")

	binID := primitive.Binary{Subtype: 4, Data: []byte(recordID)}
	filter := bson.M{"_id": binID, "user_id": userID}

	result, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("record not found or already deleted")
	}
	return nil
}

func SearchUserRecords(userID int, searchQuery string) ([]RecordAfterResponse, error) {
	coll := Client.Database("chatgpt_telegram_bot").Collection("record")

	// Выполняем поиск по ключевым словам в поле "record"
	filter := bson.M{
		"user_id": userID,
		"record": bson.M{
			"$regex": primitive.Regex{
				Pattern: searchQuery,
				Options: "i", // Опция "i" делает поиск нечувствительным к регистру
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching user records: %v", err)
	}
	defer cursor.Close(ctx)

	var records []RecordAfter
	for cursor.Next(ctx) {
		var rec RecordAfter
		if err := cursor.Decode(&rec); err != nil {
			return nil, fmt.Errorf("error decoding record: %v", err)
		}
		records = append(records, rec)
	}

	var response []RecordAfterResponse
	for _, r := range records {
		rr, err := ConvertToRecordAfterResponse(r)
		if err != nil {
			return nil, err
		}
		response = append(response, rr)
	}

	if len(response) == 0 {
		return nil, nil
	}
	return response, nil
}
