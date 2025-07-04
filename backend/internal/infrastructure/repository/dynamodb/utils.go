package dynamodb

import (
	"strconv"
)

// Helper functions for type conversion used across DynamoDB repositories

// parseInt converts string to int
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseFloat64 converts string to float64
func parseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}