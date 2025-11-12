package repo

import (
	"gorm.io/gorm"
	"github.com/Bilal-Cplusoft/sunready/internal/models"
)


type HardwareRepo struct {
	db *gorm.DB
}

func NewHardwareRepo(db *gorm.DB) *HardwareRepo {
	return &HardwareRepo{db: db}
}


func (r *HardwareRepo) ListPanels() ([]*models.Panel, error) {
	var panels []*models.Panel
	if err := r.db.Find(&panels).Error; err != nil {
		return nil, err
	}
	return panels, nil
}


func (r *HardwareRepo) ListInverters() ([]*models.Inverter, error) {
	var inverters []*models.Inverter
	if err := r.db.Find(&inverters).Error; err != nil {
		return nil, err
	}
	return inverters, nil
}

func (r *HardwareRepo) ListStorages() ([]*models.Storage, error) {
	var storages []*models.Storage
	if err := r.db.Find(&storages).Error; err != nil {
		return nil, err
	}
	return storages, nil
}


func (r *HardwareRepo) CreatePanel(panel *models.Panel) error {
	if err := r.db.Create(panel).Error; err != nil {
		return err
	}
	return nil
}


func (r *HardwareRepo) CreateInverter(inverter *models.Inverter) error {
	if err := r.db.Create(inverter).Error; err != nil {
		return err
	}
	return nil
}


func (r *HardwareRepo) CreateStorage(storage *models.Storage) error {
	if err := r.db.Create(storage).Error; err != nil {
		return err
	}
	return nil
}
