package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/stonik02/proxy_service/internal/config"
	"github.com/stonik02/proxy_service/pkg/logging"

)

type repository struct {
	logger *logging.Logger
	cfg    config.Config
}

func NewRepository(logger *logging.Logger, cfg config.Config) Repository {
	return &repository{
		logger: logger,
		cfg:    cfg,
	}
}

func (r *repository) TokenVrification(token string) PersonDataInToken {
	// Заглушка
	return PersonDataInToken{
		Id: "43fb1c66-b6ed-4af0-9906-6bbadf91aee0",
	}
}

type CustomerInfo struct {
	Id string
}

type CustomClaimsExample struct {
	*jwt.StandardClaims
	TokenType string
	CustomerInfo
}

func (r *repository) CreateJWTAccessToken(userID string) (string, error) {
	key := r.cfg.JWT.AccessKey

	payload := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(key))
	if err != nil {
		r.logger.Error(err)
		return "", err
	}
	fmt.Printf("token = %s", t)
	return t, nil
}

func (r *repository) CreateJWTRefreshToken(userID string) (string, error) {
	panic("CreateJWTRefreshToken")
}
