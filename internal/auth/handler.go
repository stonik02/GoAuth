package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stonik02/proxy_service/internal/handlers"
	"github.com/stonik02/proxy_service/internal/persons"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/middleware"
)

var _ handlers.Handler = &handler{}

const (
	registerURL = "/register"
	authURL     = "/auth"
	refreshURL  = "/refresh"
)

type handler struct {
	repository                Repository
	logger                    logging.Logger
	checkPermissionMiddleware middleware.AuthorizedRoleMiddleware
}

func NewHandler(logger logging.Logger, repository Repository, checkPermissionMiddleware middleware.AuthorizedRoleMiddleware) handlers.Handler {
	return &handler{logger: logger,
		repository:                repository,
		checkPermissionMiddleware: checkPermissionMiddleware,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.POST(registerURL, h.Registration)
	router.POST(authURL, h.checkPermissionMiddleware.BasicAuth(h.Auth, "is_admin"))
	router.POST(refreshURL, h.Refresh)

}

func (h *handler) Registration(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto RegisterDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return
	}

	newPerson, err := h.repository.RegisterPerson(context.TODO(), dto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User exist"))
		return
	}

	allBytes, err := json.Marshal(newPerson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(allBytes))
}

func (h *handler) Auth(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto person.AuthDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return
	}

	validateData, err := h.repository.Auth(context.TODO(), dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong data"))
		return
	}

	allBytes, err := json.Marshal(validateData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		return
	}

	fmt.Print(allBytes)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte([]byte("Типо токен")))
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
