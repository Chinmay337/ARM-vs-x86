package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

func main() {
	lambda.Start(HandleRequest)

}

type Subscriber struct {
	Email  string `json:"email"`
	Origin string `json:"origin_article"`
	Date   string `json:"date"`
}

type PutSubscriber struct {
	Email  string
	Origin string
	Date   string
}

func AddSubscriber(user *Subscriber) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	tableName := os.Getenv("TABLE_NAME")

	item := PutSubscriber{
		Email:  user.Email,
		Origin: user.Origin,
		Date:   user.Date,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling user: %s", err)
	}

	fmt.Println("in putitem", av, user)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("Successfully added '" + user.Email + "to subscriber list")

	return "Successfully added '" + user.Email + "to subscriber list", nil
	// snippet-end:[dynamodb.go.create_item.call]
}

func HandleRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	req, _ := json.Marshal(request)
	fmt.Println(string(req))

	ApiResponse := events.LambdaFunctionURLResponse{Body: "Success testing request!!", StatusCode: 200}
	return ApiResponse, nil

}
