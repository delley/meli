package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

// Stats is the representation of the amount of simian DNA, amount of human DNA, and the proportion of simians to humans.
type Stats struct {
	CountMutantDna int64   `json:"count_mutant_dna"`
	CountHumanDna  int64   `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

func stats(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	countMutantDna, err := countDna(true)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	countHumanDna, err := countDna(false)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	ratio := 0.0

	if countMutantDna > 0 && countHumanDna > 0 {
		ratio = float64(countMutantDna) / float64(countHumanDna)
	}

	s := Stats{
		CountMutantDna: countMutantDna,
		CountHumanDna:  countHumanDna,
		Ratio:          ratio,
	}

	response, err := json.Marshal(s)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error while decoding to string value: %s", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil

}

func countDna(isSimian bool) (int64, error) {
	// ==========================================
	// mover esse trecho para o pkg db para evitar duplicacao de codigo
	//
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return 0, fmt.Errorf("Error while retrieving AWS credentials: %s", err.Error())
	}

	svc := dynamodb.New(cfg)

	filt := expression.Name("IsSimian").Equal(expression.Value(isSimian))

	proj := expression.NamesList(expression.Name("ID"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		return 0, fmt.Errorf("Got error building expression: %s", err.Error())
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(os.Getenv("TABLE_NAME")),
	}

	// Make the DynamoDB Query API call
	req := svc.ScanRequest(params)
	res, err := req.Send()
	if err != nil {
		return 0, fmt.Errorf("Error while scanning DynamoDB: %s", err.Error())
	}

	return *res.Count, nil
}

func main() {
	lambda.Start(stats)
}
