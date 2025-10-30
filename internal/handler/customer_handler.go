package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/service"
	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}


type CreateCustomerRequest struct {
	FirstName              string   `json:"first_name" example:"John"`
	LastName               string   `json:"last_name" example:"Smith"`
	Email                  string   `json:"email" example:"john.smith@email.com"`
	PhoneNumber            *string  `json:"phone_number" example:"+1-555-123-4567"`
	Address                string   `json:"address" example:"123 Main St, San Francisco, CA 94102"`
	City                   *string  `json:"city" example:"San Francisco"`
	State                  *string  `json:"state" example:"CA"`
	ZipCode                *string  `json:"zip_code" example:"94102"`
	PropertyType           *string  `json:"property_type" example:"single_family"`
	RoofType               *string  `json:"roof_type" example:"asphalt_shingle"`
	HomeOwnershipType      *string  `json:"home_ownership_type" example:"owner"`
	AverageMonthlyBill     *float64 `json:"average_monthly_bill" example:"150.00"`
	UtilityProvider        *string  `json:"utility_provider" example:"PG&E"`
	LeadId                 *string  `json:"lead_id" example:"12345"`
	ReferralCode           *string  `json:"referral_code" example:"FRIEND2024"`
	Notes                  *string  `json:"notes" example:"Interested in 10kW system"`
	PreferredContactMethod *string  `json:"preferred_contact_method" example:"email"`
}

type UpdateCustomerStatusRequest struct {
	Status string `json:"status" example:"qualified"`
}


type UpdateCustomerRequest struct {
	FirstName              *string  `json:"first_name" example:"John"`
	LastName               *string  `json:"last_name" example:"Smith"`
	PhoneNumber            *string  `json:"phone_number" example:"+1-555-123-4567"`
	Address                *string  `json:"address" example:"123 Main St, San Francisco, CA 94102"`
	City                   *string  `json:"city" example:"San Francisco"`
	State                  *string  `json:"state" example:"CA"`
	ZipCode                *string  `json:"zip_code" example:"94102"`
	PropertyType           *string  `json:"property_type" example:"single_family"`
	RoofType               *string  `json:"roof_type" example:"asphalt_shingle"`
	HomeOwnershipType      *string  `json:"home_ownership_type" example:"owner"`
	AverageMonthlyBill     *float64 `json:"average_monthly_bill" example:"150.00"`
	UtilityProvider        *string  `json:"utility_provider" example:"PG&E"`
	Notes                  *string  `json:"notes" example:"Interested in 10kW system"`
	PreferredContactMethod *string  `json:"preferred_contact_method" example:"email"`
	Status                 *string  `json:"status" example:"prospect"`
	IsActive               *bool    `json:"is_active" example:"true"`
}

// CreateCustomer godoc
// @Summary      Create a new customer
// @Description  Creates a new customer in the system
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer  body      CreateCustomerRequest  true  "Customer data"
// @Success      201       {object}  models.Customer
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/customers [post]
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	leadId , _ := strconv.Atoi(*req.LeadId)

	customer := &models.Customer{
		FirstName:              req.FirstName,
		LastName:               req.LastName,
		Email:                  req.Email,
		PhoneNumber:            req.PhoneNumber,
		Address:                req.Address,
		City:                   req.City,
		State:                  req.State,
		ZipCode:                req.ZipCode,
		PropertyType:           req.PropertyType,
		RoofType:               req.RoofType,
		HomeOwnershipType:      req.HomeOwnershipType,
		AverageMonthlyBill:     req.AverageMonthlyBill,
		UtilityProvider:        req.UtilityProvider,
		LeadId:             leadId,
		ReferralCode:           req.ReferralCode,
		Notes:                  req.Notes,
		PreferredContactMethod: req.PreferredContactMethod,
		IsActive:               true,
	}

	if err := h.customerService.CreateCustomer(r.Context(), customer); err != nil {
		http.Error(w, "failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetCustomer godoc
// @Summary      Get customer by ID
// @Description  Retrieves a customer by their ID
// @Tags         customers
// @Produce      json
// @Param        id   path      int  true  "Customer ID"
// @Success      200  {object}  models.Customer
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/customers/{id} [get]
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
		return
	}

	customer, err := h.customerService.GetCustomerByID(r.Context(), id)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomer godoc
// @Summary      Update customer
// @Description  Updates an existing customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id        path      int                    true  "Customer ID"
// @Param        customer  body      UpdateCustomerRequest  true  "Updated customer data"
// @Success      200       {object}  models.Customer
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
		return
	}

	customer, err := h.customerService.GetCustomerByID(r.Context(), id)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.FirstName != nil {
		customer.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		customer.LastName = *req.LastName
	}
	if req.PhoneNumber != nil {
		customer.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.City != nil {
		customer.City = req.City
	}
	if req.State != nil {
		customer.State = req.State
	}
	if req.ZipCode != nil {
		customer.ZipCode = req.ZipCode
	}
	if req.PropertyType != nil {
		customer.PropertyType = req.PropertyType
	}
	if req.RoofType != nil {
		customer.RoofType = req.RoofType
	}
	if req.HomeOwnershipType != nil {
		customer.HomeOwnershipType = req.HomeOwnershipType
	}
	if req.AverageMonthlyBill != nil {
		customer.AverageMonthlyBill = req.AverageMonthlyBill
	}
	if req.UtilityProvider != nil {
		customer.UtilityProvider = req.UtilityProvider
	}
	if req.Notes != nil {
		customer.Notes = req.Notes
	}
	if req.PreferredContactMethod != nil {
		customer.PreferredContactMethod = req.PreferredContactMethod
	}
	if req.Status != nil {
		customer.Status = *req.Status
	}
	if req.IsActive != nil {
		customer.IsActive = *req.IsActive
	}

	if err := h.customerService.UpdateCustomer(r.Context(), customer); err != nil {
		http.Error(w, "failed to update customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomer godoc
// @Summary      Delete customer
// @Description  Deletes a customer by ID
// @Tags         customers
// @Param        id   path      int  true  "Customer ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
		return
	}

	if err := h.customerService.DeleteCustomer(r.Context(), id); err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCustomers godoc
// @Summary      List customers
// @Description  Retrieves a list of customers with optional pagination and filtering
// @Tags         customers
// @Produce      json
// @Param        limit   query     int     false  "Limit number of results"  default(50)
// @Param        offset  query     int     false  "Offset for pagination"   default(0)
// @Param        status  query     string  false  "Filter by status"
// @Param        search  query     string  false  "Search customers by name, email, or address"
// @Success      200     {array}   models.Customer
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/customers [get]
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	var customers []*models.Customer
	var err error

	if search != "" {
		customers, err = h.customerService.SearchCustomers(r.Context(), search, limit, offset)
	} else if status != "" {
		customers, err = h.customerService.ListCustomersByStatus(r.Context(), status, limit, offset)
	} else {
		customers, err = h.customerService.ListCustomers(r.Context(), limit, offset)
	}

	if err != nil {
		http.Error(w, "failed to list customers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// UpdateCustomerStatus godoc
// @Summary      Update customer status
// @Description  Updates the status of a customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id      path      int                    true  "Customer ID"
// @Param        status  body      UpdateCustomerStatusRequest      true  "Status update"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /api/customers/{id}/status [patch]
func (h *CustomerHandler) UpdateCustomerStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid customer id", http.StatusBadRequest)
		return
	}

	var req UpdateCustomerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	if err := h.customerService.UpdateCustomerStatus(r.Context(), id, req.Status); err != nil {
		http.Error(w, "failed to update customer status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Customer status updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCustomerStats godoc
// @Summary      Get customer statistics
// @Description  Retrieves customer statistics by status
// @Tags         customers
// @Produce      json
// @Success      200  {object}  service.CustomerStats
// @Failure      500  {object}  map[string]string
// @Router       /api/customers/stats [get]
func (h *CustomerHandler) GetCustomerStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.customerService.GetCustomerStats(r.Context())
	if err != nil {
		http.Error(w, "failed to get customer stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
