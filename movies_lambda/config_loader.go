package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/mock"
)

type ConfigLoader interface {
	LoadDefaultConfig(ctx context.Context) (cfg aws.Config, err error)
}

type MockedConfigLoader struct {
	mock.Mock
}

func NewMockedConfigLoader() *MockedConfigLoader {
	return &MockedConfigLoader{}
}

func (m *MockedConfigLoader) LoadDefaultConfig(ctx context.Context) (cfg aws.Config, err error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return aws.Config{}, args.Error(1)
	}

	return args.Get(0).(aws.Config), args.Error(1)
}

type AwsConfigLoader struct{}

func NewAwsConfigLoader() *AwsConfigLoader {
	return &AwsConfigLoader{}
}

func (l *AwsConfigLoader) LoadDefaultConfig(ctx context.Context) (cfg aws.Config, err error) {
	return config.LoadDefaultConfig(ctx)
}
