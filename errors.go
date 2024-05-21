// errors.go
package main

import "fmt"

// Custom error types
type ConfigError struct {
	Msg string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("Config Error: %s", e.Msg)
}

type DatabaseError struct {
	Msg string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Database Error: %s", e.Msg)
}

type NetworkError struct {
	Msg string
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network Error: %s", e.Msg)
}
