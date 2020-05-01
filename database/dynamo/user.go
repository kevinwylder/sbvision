package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kevinwylder/sbvision"
)

// GetUser gets the user by their email
func (sb *SBDatabase) GetUser(email string) (*sbvision.User, error) {
	output, err := sb.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(userTableName),
		Key: mustMarshalMap(map[string]string{
			"email": email,
		}, "GetUser"),
	})
	if err != nil {
		return nil, err
	}
	var user sbvision.User
	err = dynamodbattribute.UnmarshalMap(output.Item, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AddUser puts the user in the dynamodb table
func (sb *SBDatabase) AddUser(user *sbvision.User) error {
	_, err := sb.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(userTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(user.Email),
			},
			"username": {
				S: aws.String(user.Username),
			},
			"videos": {
				L: []*dynamodb.AttributeValue{},
			},
			"clips": {
				L: []*dynamodb.AttributeValue{},
			},
		},
	})
	return err
}
