package person

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/stonik02/proxy_service/internal/handlers"
	"github.com/stonik02/proxy_service/pkg/logging"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	repository Repository
	logger     logging.Logger
}

func NewHandler(logger logging.Logger, repository Repository) handlers.Handler {
	return &handler{logger: logger, repository: repository}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.GET(userURL, h.GetUserById)
	router.POST(usersURL, h.CreateUser)
	router.PATCH(userURL, h.UpdateUser)

}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	persons, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Sql error"))
	}

	allBytes, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Sql error"))
	}

	w.WriteHeader(200)
	w.Write([]byte(allBytes))
}

func (h *handler) GetUserById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("uuid")

	u, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Errorf("Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))
		return
	}

	allBytes, err := json.Marshal(u)
	if err != nil {
		h.logger.Errorf("Error: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(allBytes))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var prs Person
	err := json.NewDecoder(r.Body).Decode(&prs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}

	err = h.repository.Create(context.TODO(), &prs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}

	allBytes, err := json.Marshal(prs)

	w.WriteHeader(201)
	w.Write([]byte(allBytes))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var prs Person
	err := json.NewDecoder(r.Body).Decode(&prs)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}
	id := params.ByName("uuid")
	prs.Id = id

	err = h.repository.Update(context.TODO(), &prs)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}

	allByte, err := json.Marshal(prs)
	w.WriteHeader(201)
	w.Write([]byte(allByte))
}
