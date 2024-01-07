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

func (r *repository) TokenVrification(tokenString string, key string) (PersonDataInToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			newErr := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			r.logger.Errorf(newErr.Error())
			return nil, newErr
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(key), nil
	})
	if err != nil {
		newErr := fmt.Errorf("Error parse token")
		return PersonDataInToken{}, newErr
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return PersonDataInToken{
			Id:    fmt.Sprintf("%s", claims["uid"]),
			Email: fmt.Sprintf("%s", claims["email"]),
		}, nil
	} else {
		return PersonDataInToken{}, err
	}
}

// CreateJWTAccessToken создает jwt access token, в который включен userId.
// Токен действителен 1 час
func (r *repository) CreateJWTAccessToken(person PersonDataInToken) (string, error) {
	secretKey := r.cfg.JWT.AccessKey
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = person.Id
	claims["email"] = person.Email
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CreateJWTRefreshToken создает jwt refresh token, в который включен userId.
// Токен действителен 1 месяц
func (r *repository) CreateJWTRefreshToken(person PersonDataInToken) (string, error) {
	secretKey := r.cfg.JWT.AccessKey
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = person.Id
	claims["email"] = person.Email
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
