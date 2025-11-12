package handler

import (
	"net/http"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"encoding/json"
	"github.com/Bilal-Cplusoft/sunready/internal/models"

)

type HardwareHandler struct {
	hardwareRepo *repo.HardwareRepo
}

func NewHardwareHandler(hardwareRepo *repo.HardwareRepo) *HardwareHandler {
	return &HardwareHandler{hardwareRepo: hardwareRepo}
}


// ListPanels godoc
// @Summary List all panels
// @Description Get a list of all panels
// @Tags Hardware
// @Produce json
// @Success 200 {array} models.Panel
// @Failure 500 {object} map[string]string
// @Router /api/hardware/panels [get]
func (h *HardwareHandler) ListPanels(w http.ResponseWriter, r *http.Request) {
	panels, err := h.hardwareRepo.ListPanels()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondJSON(w, http.StatusOK, panels)
}

// ListStorages godoc
// @Summary List all storages
// @Description Get a list of all storage units
// @Tags Hardware
// @Produce json
// @Success 200 {array} models.Storage
// @Failure 500 {object} map[string]string
// @Router /api/hardware/storages [get]
func (h *HardwareHandler) ListStorages(w http.ResponseWriter, r *http.Request) {
	storages, err := h.hardwareRepo.ListStorages()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondJSON(w, http.StatusOK, storages)
}

// ListInverters godoc
// @Summary List all inverters
// @Description Get a list of all inverters
// @Tags Hardware
// @Produce json
// @Success 200 {array} models.Inverter
// @Failure 500 {object} map[string]string
// @Router /api/hardware/inverters [get]
func (h *HardwareHandler) ListInverters(w http.ResponseWriter, r *http.Request) {
	inverters, err := h.hardwareRepo.ListInverters()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondJSON(w, http.StatusOK, inverters)
}


// AddPanel godoc
// @Summary Add a new panel
// @Description Add a new panel record
// @Tags Hardware
// @Accept json
// @Produce json
// @Param panel body models.Panel true "Panel payload"
// @Success 201 {object} models.Panel
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hardware/panels [post]
func (h *HardwareHandler) AddPanel(w http.ResponseWriter, r *http.Request) {
	var panel models.Panel
	if err := json.NewDecoder(r.Body).Decode(&panel); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.hardwareRepo.CreatePanel(&panel); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create panel")
		return
	}

	respondJSON(w, http.StatusCreated, panel)
}


// AddInverter godoc
// @Summary Add a new inverter
// @Description Add a new inverter record
// @Tags Hardware
// @Accept json
// @Produce json
// @Param inverter body models.Inverter true "Inverter payload"
// @Success 201 {object} models.Inverter
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hardware/inverters [post]
func (h *HardwareHandler) AddInverter(w http.ResponseWriter, r *http.Request) {
	var inverter models.Inverter
	if err := json.NewDecoder(r.Body).Decode(&inverter); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.hardwareRepo.CreateInverter(&inverter); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create inverter")
		return
	}

	respondJSON(w, http.StatusCreated, inverter)
}


// AddStorage godoc
// @Summary Add a new storage unit
// @Description Add a new storage record
// @Tags Hardware
// @Accept json
// @Produce json
// @Param storage body models.Storage true "Storage payload"
// @Success 201 {object} models.Storage
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hardware/storages [post]
func (h *HardwareHandler) AddStorage(w http.ResponseWriter, r *http.Request) {
	var storage models.Storage
	if err := json.NewDecoder(r.Body).Decode(&storage); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.hardwareRepo.CreateStorage(&storage); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create storage")
		return
	}

	respondJSON(w, http.StatusCreated, storage)
}
