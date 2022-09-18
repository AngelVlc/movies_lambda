package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/mock"
)

type Searcher interface {
	Execute(ctx context.Context, title string) ([]Movie, error)
}

type MockedSearcher struct {
	mock.Mock
}

func NewMockedSearcher() *MockedSearcher {
	return &MockedSearcher{}
}

func (m *MockedSearcher) Execute(ctx context.Context, title string) ([]Movie, error) {
	args := m.Called(ctx, title)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Movie), args.Error(1)
}

type DynamoDbSearcher struct {
	dynamoDbApi DynamoDbApi
}

func NewDynamoDbSearcher() *DynamoDbSearcher {
	return &DynamoDbSearcher{
		dynamoDbApi: InitDynamoDbApi(),
	}
}

func (s *DynamoDbSearcher) Execute(ctx context.Context, title string) ([]Movie, error) {
	scanInput := s.getScanInput(title)

	paginator, err := s.dynamoDbApi.NewScanPaginator(ctx, scanInput)
	if err != nil {
		return nil, fmt.Errorf("error creating paginator: %v", err)
	}

	var result []Movie

	for paginator.HasMorePages() {
		scanOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error executing paginator: %v", err)
		}

		var pageMovies []Movie

		err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &pageMovies)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling the scan output items: %v", err)
		}

		result = append(result, pageMovies...)
	}

	return result, nil
}

func (s *DynamoDbSearcher) getScanInput(title string) *dynamodb.ScanInput {
	// aws dynamodb scan \
	//   --table-name Movies \
	//   --filter-expression "contains(#Title, :Title)" \
	//   --expression-attribute-names '{"#Title": "Title"}' \
	//   --expression-attribute-values '{":Title":{"S":"alien"}}'

	return &dynamodb.ScanInput{
		TableName:        aws.String("Movies"),
		FilterExpression: aws.String("contains(#Title, :Title)"),
		ExpressionAttributeNames: map[string]string{
			"#Title": "TitleToSearch",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Title": &types.AttributeValueMemberS{Value: strings.ToLower(title)},
		},
	}
}
