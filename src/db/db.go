package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// GetService create DynamoDB client sevice
func GetService() (*dynamodb.DynamoDB, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving AWS credentials: %s", err.Error())
	}

	return dynamodb.New(cfg), nil
}
