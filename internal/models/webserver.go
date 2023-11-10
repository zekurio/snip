package models

var Ok = &Status{Code: 200}

type Status struct {
	Code    int
	Message string
}
