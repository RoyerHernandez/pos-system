package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/go-chi/chi/v5"
)

type clientService interface {
	GetAll() ([]domain.Client, error)
	GetByID(id int) (*domain.Client, error)
	Create(nombre string, documento, email, telefono, direccion, fechaNacimiento *string) (int64, error)
	Update(id int, nombre string, documento, email, telefono, direccion, fechaNacimiento *string) error
	Delete(id int) error
}

// ClientController handles client CRUD endpoints.
type ClientController struct {
	service clientService
}

// NewClientController creates a new ClientController.
func NewClientController(svc clientService) (*ClientController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &ClientController{service: svc}, nil
}

// GetAll returns all clients.
func (c *ClientController) GetAll(w http.ResponseWriter, r *http.Request) {
	clients, err := c.service.GetAll()
	if err != nil {
		log.Printf("controller - GetAllClients: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch clients")
		return
	}
	RespondJSON(w, http.StatusOK, clientListFromDomain(clients))
}

// GetByID returns a single client by ID.
func (c *ClientController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid client id")
		return
	}

	client, err := c.service.GetByID(id)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - GetClientByID: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch client")
		return
	}
	if client == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "client not found")
		return
	}

	RespondJSON(w, http.StatusOK, clientResponseFromDomain(*client))
}

// Create creates a new client.
func (c *ClientController) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateClientDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	id, err := c.service.Create(req.Nombre, req.Documento, req.Email, req.Telefono, req.Direccion, req.FechaNacimiento)
	if err != nil {
		var validation apperror.ValidationError
		if errors.As(err, &validation) {
			RespondError(w, http.StatusBadRequest, "BAD_REQUEST", validation.Message)
			return
		}
		var conflict apperror.ConflictError
		if errors.As(err, &conflict) {
			RespondError(w, http.StatusConflict, "CONFLICT", conflict.Message)
			return
		}
		log.Printf("controller - CreateClient: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

// Update updates an existing client.
func (c *ClientController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid client id")
		return
	}

	var req UpdateClientDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if err := c.service.Update(id, req.Nombre, req.Documento, req.Email, req.Telefono, req.Direccion, req.FechaNacimiento); err != nil {
		var validation apperror.ValidationError
		if errors.As(err, &validation) {
			RespondError(w, http.StatusBadRequest, "BAD_REQUEST", validation.Message)
			return
		}
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - UpdateClient: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "client updated"})
}

// Delete removes a client by ID.
func (c *ClientController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid client id")
		return
	}

	if err := c.service.Delete(id); err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - DeleteClient: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete client")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "client deleted"})
}
