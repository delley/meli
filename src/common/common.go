package common

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// BuildServerError build a server error
func BuildServerError(err error) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// BuildClientError build a client error
func BuildClientError(status int, err error) (events.APIGatewayProxyResponse, error) {
	msg := http.StatusText(status)
	if err != nil {
		msg = err.Error()
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       msg,
	}, nil
}
