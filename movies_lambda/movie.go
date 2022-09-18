package main

type Movie struct {
	Title    string `dynamodbav:"Title"`
	Location string `dynamodbav:"Location"`
	Kind     string `dynamodbav:"Type"`
}
