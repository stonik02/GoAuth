package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/stonik02/proxy_service/internal/handlers"
	"github.com/stonik02/proxy_service/pkg/logging"
)

var _ handlers.Handler = &handler{}

const (
	registerURL = "/register"
	authURL     = "/auth"
	refreshURL  = "/refresh"
)

type handler struct {
	repository Repository
	logger     logging.Logger
}

func NewHandler(logger logging.Logger, repository Repository) handlers.Handler {
	return &handler{logger: logger, repository: repository}
}

func (h *handler) Register(router *httprouter.Router) {
	router.POST(registerURL, h.Registration)
	router.POST(authURL, h.Auth)
	router.POST(refreshURL, h.Refresh)

}

func (h *handler) Registration(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto RegisterDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return
	}
	// TODO:
	fmt.Printf("DTO == %s", dto)

	newPerson, err := h.repository.Register(context.TODO(), dto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User exist"))
		return
	}

	allBytes, err := json.Marshal(newPerson)
	w.WriteHeader(201)
	w.Write([]byte(allBytes))
}

func (h *handler) Auth(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}
