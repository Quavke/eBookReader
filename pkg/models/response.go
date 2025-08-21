package models

type APIResponse[T any] struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Data    T      `json:"data,omitempty"`
}