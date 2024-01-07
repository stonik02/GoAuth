package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/stonik02/proxy_service/internal/config"
	"github.com/stonik02/proxy_service/internal/token"
	utils "github.com/stonik02/proxy_service/internal/util/middleware"
	"github.com/stonik02/proxy_service/pkg/logging"
)

type AuthorizedRoleMiddleware struct {
	UtilsRepository utils.Repository
	TokenRepository token.Repository
	Cfg             config.Config
	Logger          *logging.Logger
}

func NewAuthorizedRole(utilsRepository utils.Repository, tokenRepository token.Repository, cfg config.Config, logger *logging.Logger) *AuthorizedRoleMiddleware {
	return &AuthorizedRoleMiddleware{
		UtilsRepository: utilsRepository,
		TokenRepository: tokenRepository,
		Cfg:             cfg,
		Logger:          logger,
	}
}

func (au *AuthorizedRoleMiddleware) splitBearerToken(bearerToken string) (string, error) {
	var split = strings.SplitAfterN(bearerToken, " ", 2)
	if len(split) == 1 {
		return "", fmt.Errorf("Split token error")
	}

	if split[0] != "Bearer" {
		return "", fmt.Errorf("No bearer")
	}

	return split[1], nil
}

func (au *AuthorizedRoleMiddleware) ParsePersonDataFromAccessToken(bearerToken string, accessKey string) (token.PersonDataInToken, error) {
	personData := token.PersonDataInToken{}
	token, err := au.splitBearerToken(bearerToken)
	if err != nil {
		au.Logger.Errorf("ParsePersonDataFromAccessToken: splitBearerToken error: %s", err)
		return personData, err
	}
	personData, err = au.TokenRepository.TokenVrification(token, accessKey)
	if err != nil {
		au.Logger.Errorf("ParsePersonDataFromAccessToken: TokenVrification error: %s", err)
		return personData, err
	}
	return personData, nil
}

func (au *AuthorizedRoleMiddleware) CheckingPersonRolesWithAllowedRole(userId string, allowedRole string) (bool, error) {
	personRoles, err := au.UtilsRepository.GetUserRoleNames(context.TODO(), userId)
	if err != nil {
		au.Logger.Errorf("CheckingPersonRolesWithAllowedRole: GetUserRoleNames error: %s", err)
		return false, err
	}
	for _, role := range personRoles {
		if role == allowedRole {
			return true, nil
		}
	}
	return false, nil
}

func (au *AuthorizedRoleMiddleware) BasicAuth(handle httprouter.Handle, allowedRole string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		bearerToken := r.Header.Values("Authorization")
		if bearerToken != nil {
			accessKey := au.Cfg.JWT.AccessKey
			personData, err := au.ParsePersonDataFromAccessToken(bearerToken[0], accessKey)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Вы не авторизованы"))
				return
			}
			personHasRole, err := au.CheckingPersonRolesWithAllowedRole(personData.Id, allowedRole)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			if personHasRole {
				// TODO: добавить в request данные пользователя
				handle(w, r, ps)
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Вы не авторизованы"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Вы не авторизованы"))
		}
		// 	user, password, hasAuth := r.BasicAuth()

		// 	if hasAuth && user == requiredUser && password == requiredPassword {
		// 		// Delegate request to the given handle
		// 		h(w, r, ps)
		// 	} else {
		// 		// Request Basic Authentication otherwise
		// 		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		// 		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// 	}
	}
}
