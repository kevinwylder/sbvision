package dynamo

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kevinwylder/sbvision"
)

// AddClip adds the given clip to the clips table, and updates the user record to have it's clip id
func (sb *SBDatabase) AddClip(clip *sbvision.Clip, user *sbvision.User) error {
	var id [12]byte
	rand.Read(id[:])
	clip.ID = base64.URLEncoding.EncodeToString(id[:])
	clip.Username = user.Username
	clip.UploadedAt = time.Now().Format("2006-01-02 15:04:05")
	_, err := sb.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(clipTableName),
		Item:      mustMarshalMap(clip, "AddClip put clip"),
	})
	if err != nil {
		return err
	}

	_, err = sb.db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(userTableName),
		Key: mustMarshalMap(map[string]string{
			"email": user.Email,
		}, "AddClip update user clips"),
		UpdateExpression: aws.String("SET clips = list_append(clips, :newclip)"),
		ExpressionAttributeValues: mustMarshalMap(map[string][]string{
			":newclip": {clip.ID},
		}, "AddClip update user clips new clips list expression"),
	})
	if err != nil {
		user.Clips = append(user.Clips, clip.ID)
	}
	return err
}

func (sb *SBDatabase) GetClips(trickName string) ([]sbvision.Clip, error) {
	var filterExpression *string
	var filterValues map[string]*dynamodb.AttributeValue
	if trickName != "" {
		filterExpression = aws.String("trick = :trickName")
		filterValues = mustMarshalMap(map[string]string{
			":trickName": trickName,
		}, "GetClips trick name value search")
	}
	var clips []sbvision.Clip
	var err error
	err = sb.db.ScanPages(&dynamodb.ScanInput{
		TableName:                 aws.String(clipTableName),
		FilterExpression:          filterExpression,
		ExpressionAttributeValues: filterValues,
	}, func(page *dynamodb.ScanOutput, isDone bool) bool {
		for _, item := range page.Items {
			var clip sbvision.Clip
			err = dynamodbattribute.UnmarshalMap(item, &clip)
			if err != nil {
				return false
			}
			clips = append(clips, clip)
		}
		return true
	})

	if err != nil {
		return nil, err
	}
	return clips, nil
}

func (sb *SBDatabase) GetClipByID(id string) (*sbvision.Clip, error) {
	data, err := sb.db.GetItem(&dynamodb.GetItemInput{
		Key: mustMarshalMap(map[string]string{
			"id": id,
		}, "GetClipById"),
		TableName: aws.String(clipTableName),
	})
	if err != nil {
		return nil, err
	}
	var clip sbvision.Clip
	err = dynamodbattribute.UnmarshalMap(data.Item, &clip)
	if err != nil {
		return nil, err
	}
	return &clip, nil
}
