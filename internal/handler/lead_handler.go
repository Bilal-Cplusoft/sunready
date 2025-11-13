package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"github.com/Bilal-Cplusoft/sunready/internal/service"
	"github.com/Bilal-Cplusoft/sunready/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type LeadHandler struct {
	leadRepo          *repo.LeadRepo
	leadService *service.LeadService
	userRepo *repo.UserRepo
}

type LeadResponse struct {
	Success bool `json:"success" example:"true"`
	LeadID  int  `json:"lead_id" example:"42"`
	HouseID int  `json:"house_id" example:"123"`
}

func NewLeadHandler(leadRepo *repo.LeadRepo, leadService *service.LeadService, userRepo *repo.UserRepo) *LeadHandler {
	return &LeadHandler{
		leadRepo:          leadRepo,
		leadService: leadService,
		userRepo: userRepo,
	}
}



// CreateLead godoc
// @Summary Create a new lead
// @Description Creates a new lead and optionally initiates 3D model generation
// @Tags leads
// @Accept json
// @Produce json
// @Param request body service.CreateLead true "Lead details"
// @Success 201 {object} LeadResponse "Returns lead and house IDs"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads [post]
func (h *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
    var req service.CreateLead
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
    if !ok {
		http.Error(w, "User ID missing from context", http.StatusUnauthorized)
		return
    }
    if req.ProjectID == 0 {
		respondError(w, http.StatusBadRequest, "Project ID is required")
		return
	}
	response, err := h.leadService.CreateLead(r.Context(), req, userID, req.ProjectID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	leadResponse := LeadResponse{
		Success: response.Success,
		LeadID:  response.LeadID,
		HouseID: response.HouseID,
	}
	respondJSON(w, http.StatusCreated, leadResponse)
}

// GetLead godoc
// @Summary Get a lead by ID
// @Description Retrieves a lead by its ID
// @Tags leads
// @Produce json
// @Param id path int true "Lead ID"
// @Success 200 {object} LeadResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/leads/{id} [get]
func (h *LeadHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}
	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}
	respondJSON(w, http.StatusOK, lead)
}

// GetMeshFiles godoc
// @Summary      Get 3D mesh files for a lead
// @Description  Retrieves the 3D mesh files associated with a specific lead ID
// @Tags         Leads
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Lead ID"
// @Success      200  {object}  client.ProfilesFiles3DResponse
// @Failure      400  {object}  ErrorResponse  "Invalid lead ID"
// @Failure      404  {object}  ErrorResponse  "Lead not found"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /leads/{id}/mesh-files [get]
func (h *LeadHandler) GetMeshFiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}
	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}
	files,err := h.leadService.GetMeshFiles(r.Context(),*lead.ExternalID)
	respondJSON(w, http.StatusOK, files)
}

// ListLeads godoc
// @Summary List leads
// @Description Retrieves a paginated list of leads
// @Tags leads
// @Produce json
// @Param user_id query int false "Filter by creator ID"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /admin/leads [get]
func (h *LeadHandler) ListLeads(w http.ResponseWriter, r *http.Request) {
	var userId *int
	limit := 20
	offset := 0

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userId = &id
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}


	var leads []*models.Lead
	var total int64
	var err error

	leads, total, err = h.leadRepo.List(r.Context(),  userId, limit, offset)
	if err != nil {
		log.Printf("Failed to list leads: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list leads")
		return
	}

	result := map[string]interface{}{
		"leads": leads,
		"total": total,
		"limit": limit,
		"offset": offset,
	}

	respondJSON(w, http.StatusOK, result)
}

// UpdateLead godoc
// @Summary Update a lead
// @Description Updates an existing lead
// @Tags leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Param request body map[string]interface{} true "Lead updates"
// @Success 200 {object} LeadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads/{id} [put]
func (h *LeadHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if kwhUsage, ok := updates["kwh_usage"].(float64); ok {
		lead.KwhUsage = kwhUsage
	}
	if systemSize, ok := updates["system_size"].(float64); ok {
		lead.SystemSize = systemSize
	}
	if panelCount, ok := updates["panel_count"].(float64); ok {
		lead.PanelCount = int(panelCount)
	}
	if annualProduction, ok := updates["annual_production"].(float64); ok {
		lead.AnnualProduction = annualProduction
	}

	if err := h.leadRepo.Update(r.Context(), lead); err != nil {
		log.Printf("Failed to update lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update lead")
		return
	}

	respondJSON(w, http.StatusOK, "success")
}

// DeleteLead godoc
// @Summary Delete a lead
// @Description Deletes a lead by ID
// @Tags leads
// @Param id path int true "Lead ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/leads/{id} [delete]
func (h *LeadHandler) DeleteLead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	if err := h.leadRepo.Delete(r.Context(), id); err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to delete lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete lead")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
