package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestService_Handler_Request_With_Empty_QueryString(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	s := Service{}

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{})

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "query string is empty", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "400 - Bad Request: query string is empty")
}

func TestService_Handler_Request_With_Empty_Title_In_The_QueryString(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	s := Service{}
	qs := map[string]string{
		"title": "",
	}

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{QueryStringParameters: qs})

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "the title to search is empty", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "400 - Bad Request: the title to search is empty")
}

func TestService_Handler_Request_With_Another_QueryString(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	s := Service{}
	qs := map[string]string{
		"another": "value",
	}

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{QueryStringParameters: qs})

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "the title to search is empty", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "400 - Bad Request: the title to search is empty")
}

func TestService_Handler_Request_With_Title_In_The_QueryString_Too_Short(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	s := Service{}
	qs := map[string]string{
		"title": "a",
	}

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{QueryStringParameters: qs})

	assert.Equal(t, 400, res.StatusCode)
	assert.Equal(t, "the title to search is too short", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "400 - Bad Request: the title to search is too short")
}

func TestService_Handler_Request_With_Valid_Title_In_The_QueryString_But_Searcher_Fails(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	mockedSearcher := NewMockedSearcher()
	s := Service{searcher: mockedSearcher}
	qs := map[string]string{
		"title": "value",
	}

	mockedSearcher.On("Execute", context.TODO(), "value").Return(nil, fmt.Errorf("some error"))

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{QueryStringParameters: qs})

	assert.Equal(t, 500, res.StatusCode)
	assert.Equal(t, "Internal error", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "500 - Internal Error: error executing the search: some error")

	mockedSearcher.AssertExpectations(t)
}

func TestService_Handler_Request_With_Valid_Title_In_The_QueryString_And_Searcher_Does_Not_Fail_But_The_Encoding_To_Json_Fails(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	mockedSearcher := NewMockedSearcher()
	encoder := func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("some error")
	}
	s := Service{searcher: mockedSearcher, encoder: encoder}
	qs := map[string]string{
		"title": "value",
	}
	searchResult := []Movie{}

	mockedSearcher.On("Execute", context.TODO(), "value").Return(searchResult, nil)

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{QueryStringParameters: qs})

	assert.Equal(t, 500, res.StatusCode)
	assert.Equal(t, "Internal error", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "500 - Internal Error: error marshaling the search results: some error")

	mockedSearcher.AssertExpectations(t)
}

func TestService_Handler_Request_With_Valid_Title_In_The_QueryString_And_Searcher_Does_Not_Fail_And_The_Encoder_Does_Not_Fail(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	mockedSearcher := NewMockedSearcher()
	encoder := func(v interface{}) ([]byte, error) {
		return []byte("encoded result"), nil
	}
	s := Service{searcher: mockedSearcher, encoder: encoder}
	qs := map[string]string{
		"title": "value",
	}
	searchResult := []Movie{}

	mockedSearcher.On("Execute", context.TODO(), "value").Return(searchResult, nil)

	res, err := s.handler(context.TODO(), events.LambdaFunctionURLRequest{QueryStringParameters: qs})

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "encoded result", res.Body)
	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "200 - Ok")

	mockedSearcher.AssertExpectations(t)
}
