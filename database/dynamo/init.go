package dynamo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// SBDatabase is the database container
type SBDatabase struct {
	db *dynamodb.DynamoDB
}

const videoTableName = "Videos"
const userTableName = "Users"

// FindTables will find the tables DynamoDB tables for the amazon account
func FindTables(session *session.Session) (*SBDatabase, error) {

	rand.Seed(time.Now().Unix())

	sb := &SBDatabase{
		db: dynamodb.New(session),
	}

	// look for video table
	err := sb.ensureTable(videoTableName, aws.String("id"), aws.String("S"), nil, nil)
	if err != nil {
		return nil, err
	}

	// look for the users table
	err = sb.ensureTable(userTableName, aws.String("email"), aws.String("S"), nil, nil)
	if err != nil {
		return nil, err
	}

	return sb, nil
}

func mustMarshalMap(input interface{}, log string) map[string]*dynamodb.AttributeValue {
	data, err := dynamodbattribute.MarshalMap(input)
	if err != nil {
		panic(err)
	}
	if log == "RemoveVideo User Update remove index value" {
		fmt.Println(data)
	}
	return data
}

func (sb *SBDatabase) ensureTable(name string, primaryKey, primaryType, secondaryKey, secondaryType *string) error {
	tableName := aws.String(name)
	_, err := sb.db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: tableName,
	})

	if e, ok := err.(awserr.Error); ok && e.Code() == dynamodb.ErrCodeResourceNotFoundException {
		thruput := &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		}
		if secondaryKey == nil {
			fmt.Println("Creating", name, "table with a primary key")
			_, err = sb.db.CreateTable(&dynamodb.CreateTableInput{
				TableName:             tableName,
				ProvisionedThroughput: thruput,
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					&dynamodb.AttributeDefinition{
						AttributeName: primaryKey,
						AttributeType: primaryType,
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					&dynamodb.KeySchemaElement{
						AttributeName: primaryKey,
						KeyType:       aws.String("HASH"),
					},
				},
			})
		} else {
			fmt.Println("Creating", name, "table with a primary and secondary key")
			_, err = sb.db.CreateTable(&dynamodb.CreateTableInput{
				TableName:             tableName,
				ProvisionedThroughput: thruput,
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					&dynamodb.AttributeDefinition{
						AttributeName: primaryKey,
						AttributeType: primaryType,
					},
					&dynamodb.AttributeDefinition{
						AttributeName: secondaryKey,
						AttributeType: secondaryType,
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					&dynamodb.KeySchemaElement{
						AttributeName: primaryKey,
						KeyType:       aws.String("HASH"),
					},
					&dynamodb.KeySchemaElement{
						AttributeName: secondaryKey,
						KeyType:       aws.String("RANGE"),
					},
				},
			})
		}
	}
	return err
}
