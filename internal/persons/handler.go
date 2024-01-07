package person

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stonik02/proxy_service/internal/handlers"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/middleware"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

const (
	admin     = "role_admin"
	user      = "role_user"
	superuser = "role_superuser"
	moder     = "role_moder"
)

type handler struct {
	repository Repository
	logger     logging.Logger
	middleware middleware.AuthorizedRoleMiddleware
}

func NewHandler(logger logging.Logger, repository Repository, middleware middleware.AuthorizedRoleMiddleware) handlers.Handler {
	return &handler{logger: logger, repository: repository, middleware: middleware}
}

func (h *handler) Register(router *httprouter.Router) {
	// Ручки с ограничениями

	// router.GET(usersURL, h.GetList)
	// router.GET(userURL, h.GetUserById)
	// router.POST(usersURL, h.middleware.BasicAuth(h.CreateUser, superuser))
	// router.PATCH(userURL, h.middleware.BasicAuth(h.UpdateUser, user)) // TODO: Сделать чтобы только сам юзер мог это делать или админ ??

	router.GET(usersURL, h.GetList)
	router.GET(userURL, h.GetUserById)
	router.POST(usersURL, h.CreateUser)
	router.PATCH(userURL, h.UpdateUser)

}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	persons, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sql error"))
	}

	allBytes, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sql error"))
	}

	w.WriteHeader(http.StatusOK)
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
		return
	}

	allBytes, err := json.Marshal(prs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(allBytes))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var prs Person
	err := json.NewDecoder(r.Body).Decode(&prs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}
	id := params.ByName("uuid")
	prs.Id = id

	err = h.repository.Update(context.TODO(), &prs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	allByte, err := json.Marshal(prs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(allByte))
}
