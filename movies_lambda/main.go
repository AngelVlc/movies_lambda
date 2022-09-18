package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	s := Service{
		searcher: InitSearcher(),
		encoder:  json.Marshal,
	}

	lambda.Start(s.handler)
}
