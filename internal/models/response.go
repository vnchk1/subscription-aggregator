package models

import "errors"

var ErrNotFound = errors.New("subscription not found")

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type ListResponse struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}
