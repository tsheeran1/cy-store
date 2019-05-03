package main

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Event struct {
	Age    int `json:"age"`
	Height int `json:"height"`
	Income int `json:"income"`
}

type Record struct {
	Userid string
	Age    int
	Height int
	Income int
}

type Output interface{}

func handler(ctx context.Context, e Event) (string, error) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Event 1 : ", e)
	r := Record{
		Userid: fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]),
		Age:    e.Age,
		Height: e.Height,
		Income: e.Income,
	}

	config := &aws.Config{
		Region: aws.String("us-east-2"),
	}

	sess := session.Must(session.NewSession(config))

	dbc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(r)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal object, %v", err))
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("compare-yourself"),
	}

	fmt.Println("input", input)

	data, err := dbc.PutItem(input)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println("data", data)

	return fmt.Sprintf("%v", data), nil
}

func main() {

	lambda.Start(handler)

}
