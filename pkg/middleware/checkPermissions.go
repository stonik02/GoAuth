package middleware

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stonik02/proxy_service/internal/roles"
	"github.com/stonik02/proxy_service/internal/token"

)

type AuthorizedRoleMiddleware struct {
	RolesRepository roles.Repository
	TokenRepository token.Repository
}

func NewAuthorizedRole(rolesRepository roles.Repository, tokenRepository token.Repository) *AuthorizedRoleMiddleware {
	return &AuthorizedRoleMiddleware{RolesRepository: rolesRepository, TokenRepository: tokenRepository}
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func (ar *AuthorizedRoleMiddleware) CheckingPersonRolesWithAllowedRole(personRoles []roles.Role, allowedRole string) {

}

// func (ar *AuthorizedRoleMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	bearerToken := r.Header.Values("Authorization")

// 	if bearerToken != nil {
// 		fmt.Printf("It's Middleware!!!  bearerToken = %s", bearerToken)
// 		ar.handler.ServeHTTP(w, r)
// 	}

// else {
//     w.WriteHeader(403)
// }
// }

func (au *AuthorizedRoleMiddleware) BasicAuth(handle httprouter.Handle, allowedRole string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		bearerToken := r.Header.Values("Authorization")

		if bearerToken != nil {
			fmt.Printf("It's Middleware!!!  bearerToken = %s \n %s", bearerToken, allowedRole)
			handle(w, r, ps)
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
