//go:build wireinject

package main

import (
	"os"

	"github.com/google/wire"
)

func InitSearcher() Searcher {
	if inTestingMode() {
		return initMockedSearcher()
	} else {
		return initDynamoDbSearcher()
	}
}

func initDynamoDbSearcher() Searcher {
	wire.Build(DynamoDbSearcherSet)
	return nil
}

func initMockedSearcher() Searcher {
	wire.Build(MockedSearcherSet)
	return nil
}

func InitDynamoDbApi() DynamoDbApi {
	if inTestingMode() {
		return initMockedDybamoDbApi()
	} else {
		return initAwsDynamoDbApi()
	}
}

func initAwsDynamoDbApi() DynamoDbApi {
	wire.Build(AwsDynamoDbApiSet)
	return nil
}

func initMockedDybamoDbApi() DynamoDbApi {
	wire.Build(MockedDynamoDbApiSet)
	return nil
}

func inTestingMode() bool {
	return len(os.Getenv("TESTING")) > 0
}

var DynamoDbSearcherSet = wire.NewSet(
	NewDynamoDbSearcher,
	wire.Bind(new(Searcher), new(*DynamoDbSearcher)),
)

var MockedSearcherSet = wire.NewSet(
	NewMockedSearcher,
	wire.Bind(new(Searcher), new(*MockedSearcher)),
)

var AwsDynamoDbApiSet = wire.NewSet(
	NewAwsDynamoDbApi,
	wire.Bind(new(DynamoDbApi), new(*AwsDynamoDbApi)),
)

var MockedDynamoDbApiSet = wire.NewSet(
	NewMockedDynamoDbApi,
	wire.Bind(new(DynamoDbApi), new(*MockedDynamoDbApi)),
)
