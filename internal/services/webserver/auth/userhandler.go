package auth

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/services/util/static"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db         database.IDatabase
	jwtHandler *JWTTokenHandler
}

func NewUserHandler(ctn di.Container) *UserHandler {
	return &UserHandler{
		db:         ctn.Get(static.DiDatabase).(database.IDatabase),
		jwtHandler: ctn.Get(static.DiJWTTokenHandler).(*JWTTokenHandler),
	}
}

func (h *UserHandler) RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return h.db.AddUpdateUser(user)
}

func (h *UserHandler) LoginUser(username, password string) (string, error) {
	user, err := h.db.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token, _, err := h.jwtHandler.GenerateToken(user.UUID)
	if err != nil {
		return "", err
	}

	// check if the token has expired

	return token, nil
}
