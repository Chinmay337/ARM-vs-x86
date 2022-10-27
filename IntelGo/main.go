package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"sort"
)

func main() {
	lambda.Start(HandleRequest)

}

func HandleRequest(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	req, _ := json.Marshal(request)
	fmt.Println(string(req))

	ApiResponse := events.LambdaFunctionURLResponse{Body: "Success testing request!!", StatusCode: 200}
	return ApiResponse, nil

}

func ComputeFloat(input []float64) (events.LambdaFunctionURLResponse, error) {
	sample_arr := []float64{10.4, 11.2, 192.5, 200.145, 12.1341}
	for _, val := range sample_arr {
		fmt.Println("Curr item is", val)
		closest_items := []float64{}
		bigger := 0
		smaller := 0
		rng := 1.0
		for _, v := range sample_arr {
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
		rng = rng / float64(len(sample_arr))
		sort.Float64s(closest_items)
		fmt.Println("closest items", closest_items)
		fmt.Println(bigger, "items bigger", smaller, "items smaller")
		fmt.Println("Range of number is ", rng)

	}
	ApiResponse := events.LambdaFunctionURLResponse{Body: "Successfully Completed Computation", StatusCode: 200}
	return ApiResponse, nil
}

func ComputeInt(input []int) (events.LambdaFunctionURLResponse, error) {
	sample_arr := []int{10, 11, 192, 200, 12}
	for _, val := range sample_arr {
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
