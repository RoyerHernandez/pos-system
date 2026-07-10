package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/services"
	"github.com/RoyerHernandez/pos-system/api/internal/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch users")
		return
	}

	var responses []models.UserResponse
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	RespondJSON(w, http.StatusOK, responses)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user id")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch user")
		return
	}
	if user == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}

	RespondJSON(w, http.StatusOK, user.ToResponse())
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid form data")
		return
	}

	user := &models.User{
		Usuario: r.FormValue("usuario"),
		Nombre:  r.FormValue("nombre"),
		Perfil:  r.FormValue("perfil"),
	}
	password := r.FormValue("password")

	// Handle optional photo upload
	file, header, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()
		filename, err := utils.SaveImage(file, header, "usuarios", fmt.Sprintf("%03d", 0))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload user photo")
			return
		}
		user.Foto = &filename
	}

	id, err := h.userService.Create(user, password)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create user")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid form data")
		return
	}

	user := &models.User{
		ID:      id,
		Usuario: r.FormValue("usuario"),
		Nombre:  r.FormValue("nombre"),
		Perfil:  r.FormValue("perfil"),
	}
	password := r.FormValue("password")

	// Handle optional photo upload
	file, header, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()
		filename, err := utils.SaveImage(file, header, "usuarios", fmt.Sprintf("%03d", id))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload user photo")
			return
		}
		user.Foto = &filename
	}

	if err := h.userService.Update(user, password); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update user")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user id")
		return
	}

	if err := h.userService.Delete(id); err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete user")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}
