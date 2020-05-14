package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kevinwylder/sbvision"
)

// AddVideo creates an ID for the video and adds it to the dynamo table
func (sb *SBDatabase) AddVideo(video *sbvision.Video, user *sbvision.User) error {
	video.UploadedBy = user.Username
	video.UploaderEmail = user.Email
	_, err := sb.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(videoTableName),
		Item:      mustMarshalMap(video, "AddVideo put video"),
	})
	if err != nil {
		return err
	}
	_, err = sb.db.UpdateItem(&dynamodb.UpdateItemInput{
		Key: mustMarshalMap(map[string]string{
			"email": user.Email,
		}, "AddVideo update user key"),
		TableName:        aws.String(userTableName),
		UpdateExpression: aws.String("SET videos = list_append(videos, :newid)"),
		ExpressionAttributeValues: mustMarshalMap(map[string][]string{
			":newid": {video.ID},
		}, "AddVideo update user id list expression"),
	})
	user.Videos = append(user.Videos, video.ID)
	return err
}

// GetVideoByID will lookup the video with the given id
func (sb *SBDatabase) GetVideoByID(id string) (*sbvision.Video, error) {
	data, err := sb.db.GetItem(&dynamodb.GetItemInput{
		Key: mustMarshalMap(map[string]string{
			"id": id,
		}, "GetVideoById"),
		TableName: aws.String(videoTableName),
	})
	if err != nil {
		return nil, err
	}
	var video sbvision.Video
	err = dynamodbattribute.UnmarshalMap(data.Item, &video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// GetAllVideos returns all the videos uploaded
func (sb *SBDatabase) GetAllVideos() ([]sbvision.Video, error) {
	var videos []sbvision.Video
	err := sb.db.ScanPages(&dynamodb.ScanInput{
		TableName: aws.String(videoTableName),
	}, func(page *dynamodb.ScanOutput, isDone bool) bool {
		for _, item := range page.Items {
			var video sbvision.Video
			err := dynamodbattribute.UnmarshalMap(item, &video)
			if err != nil {
				return false
			}
			videos = append(videos, video)
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideos gets a list of all the videos for that user
func (sb *SBDatabase) GetVideos(user *sbvision.User) ([]sbvision.Video, error) {
	var err error
	user, err = sb.GetUser(user.Email)
	if err != nil {
		return nil, err
	}
	keys := map[string]*dynamodb.KeysAndAttributes{
		videoTableName: {
			Keys: []map[string]*dynamodb.AttributeValue{},
		},
	}

	for _, id := range user.Videos {
		keys[videoTableName].Keys = append(keys[videoTableName].Keys, mustMarshalMap(map[string]string{
			"id": id,
		}, "GetVideos user id"))
	}

	var videos []sbvision.Video
	err = sb.db.BatchGetItemPages(&dynamodb.BatchGetItemInput{
		RequestItems: keys,
	}, func(page *dynamodb.BatchGetItemOutput, isDone bool) bool {
		for _, response := range page.Responses[videoTableName] {
			var video sbvision.Video
			err = dynamodbattribute.UnmarshalMap(response, &video)
			if err != nil {
				return false
			}
			videos = append(videos, video)
		}
		return true
	})

	if err != nil {
		return nil, err
	}
	return videos, nil
}

// RemoveVideo removes the reference to that video from the database
// UploaderEmail is required to be set in the video struct for this method to be successful
func (sb *SBDatabase) RemoveVideo(video *sbvision.Video) error {
	_, err := sb.db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(videoTableName),
		Key: mustMarshalMap(map[string]string{
			"id": video.ID,
		}, "RemoveVideo DeleteItem"),
	})
	if err != nil {
		return err
	}
	if video.UploaderEmail == "" {
		return fmt.Errorf("Uploader Email required to delete this video")
	}

	user, err := sb.GetUser(video.UploaderEmail)
	if err != nil {
		return err
	}
	deleteID := -1
	for i := range user.Videos {
		if user.Videos[i] == video.ID {
			deleteID = i
			break
		}
	}
	if deleteID == -1 {
		return fmt.Errorf("Could not find the video to delete from user profile")
	}
	_, err = sb.db.UpdateItem(&dynamodb.UpdateItemInput{
		Key: mustMarshalMap(map[string]string{
			"email": video.UploaderEmail,
		}, "RemoveVideo User Update Key"),
		TableName:        aws.String(userTableName),
		UpdateExpression: aws.String(fmt.Sprintf("REMOVE videos[%d]", deleteID)),
	})
	if err != nil {
		return err
	}
	return err
}
