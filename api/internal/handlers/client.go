package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
	clientService *services.ClientService
}

func NewClientHandler(clientService *services.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

func (h *ClientHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	clients, err := h.clientService.GetAll()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch clients")
		return
	}
	RespondJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid client id")
		return
	}

	client, err := h.clientService.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch client")
		return
	}
	if client == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "client not found")
		return
	}

	RespondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	id, err := h.clientService.Create(req.Nombre, req.Documento, req.Email, req.Telefono, req.Direccion, req.FechaNacimiento)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid client id")
		return
	}

	var req models.UpdateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if err := h.clientService.Update(id, req.Nombre, req.Documento, req.Email, req.Telefono, req.Direccion, req.FechaNacimiento); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "client updated"})
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid client id")
		return
	}

	if err := h.clientService.Delete(id); err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete client")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "client deleted"})
}
