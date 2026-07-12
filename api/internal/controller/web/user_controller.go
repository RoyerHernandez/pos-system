package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/RoyerHernandez/pos-system/api/pkg/generic"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	GetAll() ([]domain.User, error)
	GetByID(id int) (*domain.User, error)
	Create(user *domain.User, password string) (int64, error)
	Update(user *domain.User, password string) error
	Delete(id int) error
}

// UserController handles user CRUD endpoints.
type UserController struct {
	service userService
}

// NewUserController creates a new UserController.
func NewUserController(svc userService) (*UserController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &UserController{service: svc}, nil
}

// GetAll returns all users.
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAll()
	if err != nil {
		log.Printf("controller - GetAllUsers: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch users")
		return
	}

	responses := make([]UserResponseDTO, len(users))
	for i, u := range users {
		responses[i] = userResponseFromDomain(u.ToResponse())
	}

	RespondJSON(w, http.StatusOK, responses)
}

// GetByID returns a single user by ID.
func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user id")
		return
	}

	user, err := c.service.GetByID(id)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - GetUserByID: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch user")
		return
	}
	if user == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}

	RespondJSON(w, http.StatusOK, userResponseFromDomain(user.ToResponse()))
}

// Create creates a new user from multipart form data.
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid form data")
		return
	}

	user := &domain.User{
		Usuario: r.FormValue("usuario"),
		Nombre:  r.FormValue("nombre"),
		Perfil:  r.FormValue("perfil"),
	}
	password := r.FormValue("password")

	// Handle optional photo upload
	file, header, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()
		filename, err := generic.SaveImage(file, header, "usuarios", fmt.Sprintf("%03d", 0))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload user photo")
			return
		}
		user.Foto = &filename
	}

	id, err := c.service.Create(user, password)
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
		log.Printf("controller - CreateUser: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create user")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

// Update updates an existing user from multipart form data.
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid form data")
		return
	}

	user := &domain.User{
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
		filename, err := generic.SaveImage(file, header, "usuarios", fmt.Sprintf("%03d", id))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload user photo")
			return
		}
		user.Foto = &filename
	}

	if err := c.service.Update(user, password); err != nil {
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
		log.Printf("controller - UpdateUser: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update user")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "user updated"})
}

// Delete soft-deletes a user by ID.
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user id")
		return
	}

	if err := c.service.Delete(id); err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - DeleteUser: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete user")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}
