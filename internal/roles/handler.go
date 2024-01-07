package roles

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stonik02/proxy_service/internal/handlers"
	"github.com/stonik02/proxy_service/pkg/logging"
	"github.com/stonik02/proxy_service/pkg/middleware"
)

const (
	userRoleURL = "/roles/:uuid"
	rolesURL    = "/roles"
	serverError = "Internal Server Error"
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

	// router.GET(userRoleURL, h.middleware.BasicAuth(h.UserRolesList, user))
	// router.GET(rolesURL, h.middleware.BasicAuth(h.RolesList, moder))
	// router.POST(rolesURL, h.middleware.BasicAuth(h.AssignRole, superuser))
	// router.DELETE(rolesURL, h.middleware.BasicAuth(h.TakeRole, superuser))

	router.GET(userRoleURL, h.UserRolesList)
	router.GET(rolesURL, h.RolesList)
	router.POST(rolesURL, h.AssignRole)
	router.DELETE(rolesURL, h.TakeRole)
}

func (h *handler) RolesList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	roles, err := h.repository.GetAllRoles(context.TODO())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(serverError))
		return
	}
	allBytes, _ := json.Marshal(roles)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(allBytes))
}

func (h *handler) UserRolesList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("uuid")

	response, err := h.repository.GetUserWithRoles(context.TODO(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	allBytes, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(allBytes))
}

func (h *handler) AssignRole(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto AssignRoleDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		h.logger.Errorf("Error = %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrond data"))
		return
	}
	h.logger.Tracef("DTODTODTO %s", dto)

	err = h.repository.AssignRole(context.TODO(), dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrond data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("1"))

}

func (h *handler) TakeRole(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var dto TakeRoleDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrond data"))
		return
	}
	err = h.repository.TakeRole(context.TODO(), dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrond data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("1"))
}
