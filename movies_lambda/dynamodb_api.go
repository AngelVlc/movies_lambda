package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type DynamoDbApi interface {
	NewScanPaginator(ctx context.Context, params *dynamodb.ScanInput) (Paginator, error)
}

type MockedDynamoDbApi struct {
	mock.Mock
}

func NewMockedDynamoDbApi() *MockedDynamoDbApi {
	return &MockedDynamoDbApi{}
}

func (m *MockedDynamoDbApi) NewScanPaginator(ctx context.Context, params *dynamodb.ScanInput) (Paginator, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(Paginator), args.Error(1)
}

type AwsDynamoDbApi struct{}

func NewAwsDynamoDbApi() *AwsDynamoDbApi {
	return &AwsDynamoDbApi{}
}

func (a *AwsDynamoDbApi) NewScanPaginator(ctx context.Context, params *dynamodb.ScanInput) (Paginator, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading the default config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return dynamodb.NewScanPaginator(client, params), nil
}

type Paginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type MockedPaginator struct {
	mock.Mock
}

func (m *MockedPaginator) HasMorePages() bool {
	args := m.Called()

	return args.Bool(0)
}

func (m *MockedPaginator) NextPage(ctx context.Context, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	args := m.Called(ctx, optFns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

type AwsDynamoDbPaginator struct {
	dynamodb.ScanPaginator
}
