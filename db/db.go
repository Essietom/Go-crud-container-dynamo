package db

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

func createTable() error {
    // Create a new DynamoDB client
    svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-central-1"))

    // Define the table's schema
    params := &dynamodb.CreateTableInput{
        AttributeDefinitions: []*dynamodb.AttributeDefinition{
            {
                AttributeName: aws.String("id"),
                AttributeType: aws.String("S"),
            },
        },
        KeySchema: []*dynamodb.KeySchemaElement{
            {
                AttributeName: aws.String("id"),
                KeyType:       aws.String("HASH"),
            },
        },
        TableName: aws.String("users"),
        ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        },
    }

    // Create the table
    _, err := svc.CreateTable(params)
    if err != nil {
        return err
    }

    return nil
}