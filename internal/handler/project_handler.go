package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}


// Create godoc
// @Summary Create a new project
// @Description Creates a new project record in the system.
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project data"
// @Success 201 {object} models.Project "Project created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to create project"
// @Router /projects [post]
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.projectService.Create(r.Context(), &project); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	respondJSON(w, http.StatusCreated, project)
}


// GetByID godoc
// @Summary Get project by ID
// @Description Retrieves a project by its unique ID.
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project "Project found"
// @Failure 400 {object} ErrorResponse "Invalid project ID"
// @Failure 404 {object} ErrorResponse "Project not found"
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return
	}

	respondJSON(w, http.StatusOK, project)
}



// Update godoc
// @Summary Update a project
// @Description Updates an existing project by its ID.
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Updated project data"
// @Success 200 {object} models.Project "Project updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or project ID"
// @Failure 500 {object} ErrorResponse "Failed to update project"
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project.ID = id
	if err := h.projectService.Update(r.Context(), &project); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	respondJSON(w, http.StatusOK, project)
}


// Delete godoc
// @Summary Delete a project
// @Description Deletes a project by its ID.
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 204 "Project deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid project ID"
// @Failure 500 {object} ErrorResponse "Failed to delete project"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if err := h.projectService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListByUser godoc
// @Summary List projects by user ID
// @Description Retrieves all projects associated with a specific user with pagination.
// @Tags Projects
// @Produce json
// @Param user_id query int true "User ID"
// @Param limit query int false "Results limit (default: 20)"
// @Param offset query int false "Results offset"
// @Success 200 {array} models.Project "List of projects"
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 500 {object} ErrorResponse "Failed to fetch projects"
// @Router /projects/user [get]
func (h *ProjectHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	projects, err := h.projectService.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	respondJSON(w, http.StatusOK, projects)
}
