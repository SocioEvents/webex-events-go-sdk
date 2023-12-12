package main

type Extensions struct {
	Code               string
	Cost               int
	AvailableCost      int
	Threshold          int
	DailyThreshold     int
	DailyAvailableCost int
	ReferenceId        string
	Errors             map[string]any
}

type ErrorResponse struct {
	Message       string
	Extensions    Extensions
	GraphqlErrors []map[string]any `json:"errors"`
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		GraphqlErrors: make([]map[string]any, 0),
	}
}
