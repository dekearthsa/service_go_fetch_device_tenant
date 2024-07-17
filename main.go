package main

import (
	"context"
	"encoding/json"
	"net/http"
	controller "service_go_fetch_device_tenant/controller"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var token = req.Headers["authorization"]
	status, tenan, userType, err := controller.ValidateToken(token)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       tenan,
		}, err
	}
	if status != 200 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       tenan,
		}, err
	}
	data, err := controller.HaddleFetchData(tenan, userType)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "StatusInternalServerError",
		}, err
	}

	responseBody, err := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
