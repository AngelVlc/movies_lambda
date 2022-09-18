package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type (
	jsonEncoder func(v interface{}) ([]byte, error)
)

type Service struct {
	searcher Searcher
	encoder  jsonEncoder
}

func (s *Service) handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	titleToSearch, err := s.getTitleToSearchFromRequest(request)
	if err != nil {
		log.Printf("400 - Bad Request: %v", err)

		return s.responseFor(fmt.Sprint(err), 400), nil
	}

	result, err := s.searcher.Execute(ctx, titleToSearch)
	if err != nil {
		log.Printf("500 - Internal Error: %v", fmt.Errorf("error executing the search: %v", err))

		return s.responseFor("Internal error", 500), nil
	}

	resultBytes, err := s.encoder(result)
	if err != nil {
		log.Printf("500 - Internal Error: %v", fmt.Errorf("error marshaling the search results: %v", err))

		return s.responseFor("Internal error", 500), nil
	}

	log.Printf("200 - Ok")

	return s.responseFor(string(resultBytes), 200), nil
}

func (s *Service) getTitleToSearchFromRequest(request events.LambdaFunctionURLRequest) (string, error) {
	titleToSearch := s.titleFromQueryString(request)
	log.Printf("Title to search: %q", titleToSearch)

	if len(request.QueryStringParameters) == 0 {
		return "", fmt.Errorf("query string is empty")
	}

	//TODO: check if the map has a title key

	if len(titleToSearch) == 0 {
		return "", fmt.Errorf("the title to search is empty")
	}

	if len(titleToSearch) < 3 {
		return "", fmt.Errorf("the title to search is too short")
	}

	return titleToSearch, nil
}

func (s *Service) titleFromQueryString(request events.LambdaFunctionURLRequest) string {
	return request.QueryStringParameters["title"]
}

func (s *Service) responseFor(body string, status int) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body:       body,
		StatusCode: status,
	}
}
