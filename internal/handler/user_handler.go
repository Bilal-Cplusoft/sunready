package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}


// GetByID godoc
// @Summary      Get user by ID
// @Description  Retrieves a user by their unique ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string  "Invalid user ID"
// @Failure      404  {object}  map[string]string  "User not found"
// @Router       /admin/users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}


// Update godoc
// @Summary      Update user
// @Description  Updates an existing user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int           true  "User ID"
// @Param        user  body      models.User   true  "User update payload"
// @Success      200   {object}  models.User
// @Failure      400   {object}  map[string]string  "Invalid user ID or request body"
// @Failure      500   {object}  map[string]string  "Failed to update user"
// @Router       /admin/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user.ID = id
	if err := h.userService.Update(r.Context(), &user); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}


// List godoc
// @Summary      List users
// @Description  Retrieves a paginated list of users for a specific company
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit       query     int  false  "Limit (default: 20)"
// @Param        offset      query     int  false  "Offset (default: 0)"
// @Success      200         {array}   models.User
// @Failure      400         {object}  map[string]string  "Invalid company ID"
// @Failure      500         {object}  map[string]string  "Failed to fetch users"
// @Router       /api/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, err := h.userService.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondJSON(w, http.StatusOK, users)
}


// Delete godoc
// @Summary      Delete user
// @Description  Deletes a user by their unique ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "User ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string  "Invalid user ID"
// @Failure      500  {object}  map[string]string  "Failed to delete user"
// @Router       /admin/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
