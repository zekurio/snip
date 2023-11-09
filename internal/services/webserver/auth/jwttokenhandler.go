package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/snip/internal/embedded"
	"github.com/zekurio/snip/internal/models"
	"github.com/zekurio/snip/internal/services/util/static"
)

var (
	jwtGenerationMethod = jwt.SigningMethodHS256
)

type JWTTokenHandler struct {
	sessionExpiration time.Duration
	sessionSecret     []byte
}

func NewJWTTokenHandler(ctn di.Container) *JWTTokenHandler {
	cfg := ctn.Get(static.DiConfig).(models.Config)

	return &JWTTokenHandler{
		sessionExpiration: time.Duration(cfg.Webserver.AccessToken.LifetimeSeconds) * time.Second,
		sessionSecret:     []byte(cfg.Webserver.AccessToken.Secret),
	}
}

func (h *JWTTokenHandler) GenerateToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(h.sessionExpiration)

	claims := jwt.RegisteredClaims{}
	claims.Issuer = fmt.Sprintf("snip v.%s", embedded.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.IssuedAt = jwt.NewNumericDate(now)

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(h.sessionSecret)
	return
}

func (h *JWTTokenHandler) ValidateToken(token string) (ident string, err error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(h.sessionSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		ident, ok = claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid claims")
		}
	} else {
		return "", errors.New("invalid token")
	}

	return ident, nil
}
