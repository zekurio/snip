package models

type LoginSuccessResponse struct {
	User     User
	Redirect string
}
