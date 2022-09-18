package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDynamoDbSearcher_Execute_Uses_The_Expected_ScanInput(t *testing.T) {
	mockedDynamoDbApi := &MockedDynamoDbApi{}
	searcher := &DynamoDbSearcher{mockedDynamoDbApi}
	scanInput := &dynamodb.ScanInput{
		TableName:        aws.String("Movies"),
		FilterExpression: aws.String("contains(#Title, :Title)"),
		ExpressionAttributeNames: map[string]string{
			"#Title": "TitleToSearch",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Title": &types.AttributeValueMemberS{Value: "titletosearch"},
		},
	}
	mockedPaginator := &MockedPaginator{}
	mockedPaginator.On("HasMorePages").Return(false)
	mockedDynamoDbApi.On("NewScanPaginator", context.TODO(), scanInput).Return(mockedPaginator, nil)

	res, err := searcher.Execute(context.TODO(), "TitleToSearch")

	assert.Equal(t, []Movie(nil), res)
	assert.Nil(t, err)
	mockedDynamoDbApi.AssertExpectations(t)
	mockedPaginator.AssertExpectations(t)
}

func TestDynamoDbSearcher_Execute_Returns_An_Error_When_Create_A_New_Paginator_Fails(t *testing.T) {
	mockedDynamoDbApi := &MockedDynamoDbApi{}
	searcher := &DynamoDbSearcher{mockedDynamoDbApi}
	mockedDynamoDbApi.On("NewScanPaginator", context.TODO(), mock.AnythingOfType("*dynamodb.ScanInput")).Return(nil, fmt.Errorf("some error"))

	res, err := searcher.Execute(context.TODO(), "TitleToSearch")

	assert.Nil(t, res)
	assert.EqualError(t, err, "error creating paginator: some error")
	mockedDynamoDbApi.AssertExpectations(t)
}

func TestDynamoDbSearcher_Execute_Returns_And_Error_When_PaginatorNextPage_Fails(t *testing.T) {
	mockedDynamoDbApi := &MockedDynamoDbApi{}
	searcher := &DynamoDbSearcher{mockedDynamoDbApi}
	mockedPaginator := &MockedPaginator{}
	mockedDynamoDbApi.On("NewScanPaginator", context.TODO(), mock.AnythingOfType("*dynamodb.ScanInput")).Return(mockedPaginator, nil)
	mockedPaginator.On("HasMorePages").Return(true)
	mockedPaginator.On("NextPage", context.TODO(), mock.Anything).Return(nil, fmt.Errorf("some error"))

	res, err := searcher.Execute(context.TODO(), "TitleToSearch")

	assert.Nil(t, res)
	assert.EqualError(t, err, "error executing paginator: some error")
	mockedDynamoDbApi.AssertExpectations(t)
	mockedPaginator.AssertExpectations(t)
}

func TestDynamoDbSearcher_Execute_Returns_The_Expected_Results_From_The_Paginator(t *testing.T) {
	mockedDynamoDbApi := &MockedDynamoDbApi{}
	searcher := &DynamoDbSearcher{mockedDynamoDbApi}
	mockedPaginator := &MockedPaginator{}
	scanOutput := &dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{
			{
				"Title":    &types.AttributeValueMemberS{Value: "Title 1"},
				"Location": &types.AttributeValueMemberS{Value: "loc 1"},
				"Type":     &types.AttributeValueMemberS{Value: "kind"},
			},
		},
	}
	mockedDynamoDbApi.On("NewScanPaginator", context.TODO(), mock.AnythingOfType("*dynamodb.ScanInput")).Return(mockedPaginator, nil)
	mockedPaginator.On("HasMorePages").Return(true).Once()
	mockedPaginator.On("HasMorePages").Return(false).Once()
	mockedPaginator.On("NextPage", context.TODO(), mock.Anything).Return(scanOutput, nil)

	res, err := searcher.Execute(context.TODO(), "TitleToSearch")

	assert.Equal(t, []Movie{{Title: "Title 1", Location: "loc 1", Kind: "kind"}}, res)
	assert.Nil(t, err)
	mockedDynamoDbApi.AssertExpectations(t)
	mockedPaginator.AssertExpectations(t)
}
