package auth

import (
	"errors"
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/services/database"
	"github.com/zekurio/snip/internal/util/static"
	"github.com/zekurio/snip/pkg/randutils"
)

type RefreshTokenHandlerImpl struct {
	db database.IDatabase
}

func NewRefreshTokenHandlerImpl(container di.Container) *RefreshTokenHandlerImpl {
	return &RefreshTokenHandlerImpl{
		db: container.Get(static.DiDatabase).(database.IDatabase),
	}
}

// GetRefreshToken returns a refresh token for the given ident, and saves it to the database.
func (rth *RefreshTokenHandlerImpl) GetRefreshToken(ident string) (token string, err error) {
	token, err = randutils.GetRandBase64Str(64)
	if err != nil {
		return
	}

	err = rth.db.SetUserRefreshToken(ident, token, time.Now().Add(static.AuthSessionExpiration))
	return
}

func (rth *RefreshTokenHandlerImpl) ValidateRefreshToken(token string) (ident string, err error) {
	ident, expires, err := rth.db.GetUserByRefreshToken(token)
	if err != nil {
		return
	}

	if expires.Before(time.Now()) {
		err = rth.RevokeToken(ident)
		return
	}

	user, err := rth.db.GetUserByID(ident)
	if user == nil || err != nil {
		err = errors.New("user not found")
		return
	}

	return
}

func (rth *RefreshTokenHandlerImpl) RevokeToken(ident string) error {
	err := rth.db.RevokeUserRefreshToken(ident)
	if err != nil {
		return err
	}
	return err
}
