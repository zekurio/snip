package models

import "time"

var Ok = &Status{Code: 200}

type Status struct {
	Code    int
	Message string
}

type UserLoginDetails struct {
	Username string
	Password string
}

type AccessTokenResponse struct {
	Token   string
	Expires time.Time
}
