package models

import "time"

var Ok = &Status{Code: 200}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserLoginDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type ListResponse[T any] struct {
	N    int `json:"n"`
	Data []T `json:"data"`
}

func NewListResponse[T any](data []T) ListResponse[T] {
	return ListResponse[T]{len(data), data}
}
