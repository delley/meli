package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/delley/meli/common"

	"github.com/delley/meli/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/delley/meli/domain"
)

type chain struct {
	Dna      []string `json:"dna"`
	IsSimian bool     `json:"is_simian"`
}

func process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var c chain
	err := json.Unmarshal([]byte(request.Body), &c)
	if err != nil {
		return common.BuildClientError(http.StatusBadRequest, fmt.Errorf("Invalid payload: %s", err.Error()))
	}

	isSimian, err := domain.IsSimian(c.Dna)
	if err != nil {
		return common.BuildClientError(http.StatusBadRequest, err)
	}

	c.IsSimian = isSimian

	id := buildID(c.Dna)

	svc, err := db.GetService()
	if err != nil {
		return common.BuildServerError(err)
	}

	rget := svc.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]dynamodb.AttributeValue{
			"ID": dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
	})

	res, err := rget.Send()
	if err != nil {
		return common.BuildServerError(fmt.Errorf("Error while fetching movie from DynamoDB: %s", err.Error()))
	}

	if len(res.Item) == 0 {
		rput := svc.PutItemRequest(&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("TABLE_NAME")),
			Item: map[string]dynamodb.AttributeValue{
				"ID": dynamodb.AttributeValue{
					S: aws.String(id),
				},
				"IsSimian": dynamodb.AttributeValue{
					BOOL: aws.Bool(c.IsSimian),
				},
			},
		})

		_, err = rput.Send()
		if err != nil {
			return common.BuildServerError(fmt.Errorf("Error while inserting DNA to DB: %s", err.Error()))
		}
	}

	if isSimian {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusForbidden,
	}, nil
}

func buildID(dna []string) string {
	return strings.Join(dna, "-")
}

func main() {
	lambda.Start(process)
}
