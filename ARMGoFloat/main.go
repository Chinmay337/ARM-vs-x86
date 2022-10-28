package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)

}

type InputType struct {
	Input []float64 `json:"input"`
	// Flag  string    `json:"flag"`
	// InputInt []int     `json:"input"`
}

//type Subscriber struct {
//	Email  string `json:"email"`
//	Origin string `json:"origin_article"`
//	Date   string `json:"date"`
//}
//
//type PutSubscriber struct {
//	Email  string
//	Origin string
//	Date   string
//}

//func AddSubscriber(user *Subscriber) (string, error) {
//	sess := session.Must(session.NewSessionWithOptions(session.Options{
//		SharedConfigState: session.SharedConfigEnable,
//	}))
//	svc := dynamodb.New(sess)
//
//	tableName := os.Getenv("TABLE_NAME")
//
//	item := PutSubscriber{
//		Email:  user.Email,
//		Origin: user.Origin,
//		Date:   user.Date,
//	}
//
//	av, err := dynamodbattribute.MarshalMap(item)
//	if err != nil {
//		log.Fatalf("Got error marshalling user: %s", err)
//	}
//
//	fmt.Println("in putitem", av, user)
//
//	input := &dynamodb.PutItemInput{
//		Item:      av,
//		TableName: aws.String(tableName),
//	}
//
//	_, err = svc.PutItem(input)
//	if err != nil {
//		log.Fatalf("Got error calling PutItem: %s", err)
//	}
//
//	fmt.Println("Successfully added '" + user.Email + "to subscriber list")
//
//	return "Successfully added '" + user.Email + "to subscriber list", nil
//	// snippet-end:[dynamodb.go.create_item.call]
//}

func HandleRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	//	req, _ := json.Marshal(request)
	//	fmt.Println(string(req))

	body := InputType{}
	json.Unmarshal([]byte(request.Body), &body)

	fmt.Printf("parsed", body)

	//	x, _ := json.MarshalIndent(body, "", "  ")
	//	ApiResponse := events.LambdaFunctionURLResponse{Body: string(x), StatusCode: 200}
	//	return ApiResponse, nil

	// if body.Flag == "int" {
	// 	fmt.Println("Doing Int")
	// 	return ComputeInt(body.InputInt)
	// }
	fmt.Println("Doing Float")
	return ComputeFloat(body.Input)
}

func ComputeFloat(input []float64) (events.LambdaFunctionURLResponse, error) {
	//	sample_arr := []float64{10.4, 11.2, 192.5, 200.145, 12.1341}
	for _, val := range input {
		fmt.Println("Curr item is", val)
		closest_items := []float64{}
		bigger := 0
		smaller := 0
		rng := 1.0
		for _, v := range input {
			comparator := val / v
			rng *= comparator
			if comparator > 1 {
				bigger++
			} else if comparator < 1 {
				smaller++
			}

			closest_items = append(closest_items, comparator)

			fmt.Println(comparator, val, v, "val/v")
		}
		rng = rng / float64(len(input))
		sort.Float64s(closest_items)
		fmt.Println("closest items", closest_items)
		fmt.Println(bigger, "items bigger", smaller, "items smaller")
		fmt.Println("Range of number is ", rng)

	}
	ApiResponse := events.LambdaFunctionURLResponse{Body: "Successfully Completed Computation", StatusCode: 200}
	return ApiResponse, nil
}

func ComputeInt(input []int) (events.LambdaFunctionURLResponse, error) {
	//sample_arr := []int{10, 11, 192, 200, 12}
	for _, val := range input {
		counter := 0
		for i := 2; i <= val; i++ {
			ok := isPrime(i)
			if ok {
				counter++
			}
		}
		fmt.Println("There are ", counter, "prime numbers upto ", val)

	}
	ApiResponse := events.LambdaFunctionURLResponse{Body: "Successfully Completed Computation", StatusCode: 200}
	return ApiResponse, nil
}

func isPrime(num int) bool {
	i := 2

	for i <= num {
		rem := num % i
		if rem == 0 {
			return false
		}
		if i*i > num {
			return true
		}
		i += 1

	}
	return true
}
