package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type InputType struct {
	Input []int `json:"input"`
}

func main() {
	lambda.Start(HandleRequest)

}

func HandleRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	body := InputType{}
	json.Unmarshal([]byte(request.Body), &body)

	fmt.Printf("parsed", body)
	fmt.Println("Doing Int")
	return ComputeInt(body.Input)

}

func ComputeFloat(input []float64) (events.LambdaFunctionURLResponse, error) {
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
