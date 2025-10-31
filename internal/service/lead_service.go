package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"github.com/Bilal-Cplusoft/sunready/internal/client"
)

type CreateLeadResponse struct {
	Success bool `json:"success"`
	LeadID  int  `json:"lead_id"`
	HouseID int  `json:"house_id"`
}

type LeadService struct {
	leadRepo *repo.LeadRepo
	houseRepo *repo.HouseRepo
	genabilityClient *client.Agent
}

type CreateLead struct {
	ProjectID       int     `json:"project_id" example:"1"`
	UserId        *int    `json:"user_id,omitempty" example:"1"`
	Latitude         float64 `json:"latitude" example:"37.7749"`
	Longitude        float64 `json:"longitude" example:"-122.4194"`
	SystemSize       float64 `json:"system_size" example:"10.5"`
	PanelCount       int     `json:"panel_count" example:"30"`
	HardwareType     *string `json:"hardware_type,omitempty"`
	KwhUsage         float64 `json:"kwh_usage" example:"12000"`
	Consumption      []int   `json:"consumption,omitempty"`
	LseId             int     `json:"lse_id"`
	Period            string  `json:"period"`
	TargetSolarOffset int     `json:"target_solar_offset"`
	Mode              *string `json:"mode,omitempty"`
	Unit              string  `json:"unit"`
}


func NewLeadService(leadRepo *repo.LeadRepo, houseRepo *repo.HouseRepo) *LeadService {
	var genClient *client.Agent

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Failed to initialize Genability client: %v", r)
		}
	}()

	genClient = client.NewAgent()

	return &LeadService{
		leadRepo: leadRepo,
		houseRepo: houseRepo,
		genabilityClient: genClient,
	}
}

func (s *LeadService) CreateLead(ctx context.Context, req CreateLead, userID int, projectID int) (*CreateLeadResponse, error) {
	house := models.House{
			Lat:         req.Latitude,
			Lng:         req.Longitude,
			Diameter:    20,
			Probability: 0.70,
			State:       "pending",
		}
	if err := s.houseRepo.Create(ctx, &house); err != nil {
		return nil, err
	}
	houseID := house.ID

	lead := models.Lead{
			ProjectID: projectID,
			UserID:  &userID,
			Latitude:   req.Latitude,
			Longitude:  req.Longitude,
			KwhUsage:   req.KwhUsage,
			SystemSize: req.SystemSize,
			PanelCount: req.PanelCount,
		}
		if _, err := s.leadRepo.GetLeadWithProjectByProjectID(ctx, projectID); err != nil {
			return nil, fmt.Errorf("project not found: %w", err)
        }
        if _, err := s.leadRepo.GetLeadWithUserByUserID(ctx, userID); err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
        }
		if err := s.leadRepo.Create(ctx, &lead); err != nil {
			return nil, fmt.Errorf("failed to create lead: %w", err)
		}
        Lead,err := s.leadRepo.GetLeadWithUserByUserID(ctx,userID)
        if err != nil {
            return nil, fmt.Errorf("failed to get lead with user by user id: %w", err)
        }
		if s.genabilityClient != nil {
			accountInput := client.Account{
				Address: client.AccountAddress{
					String:    *Lead.User.Address,
					Latitude:  lead.Latitude,
					Longitude: lead.Longitude,
				},
			}
			accounts := client.NewAccounts(s.genabilityClient)
			genAcc, err := accounts.Create(ctx, accountInput)
			if err != nil {
				log.Printf("Warning: Failed to create Genability account: %v", err)
			} else if genAcc != nil {
				tariffs := client.NewTariffs(s.genabilityClient)
				tariff, err := tariffs.GetCurrent(ctx, genAcc.ID)
				if err != nil {
					log.Printf("Warning: Failed to get current tariff: %v", err)
				} else if tariff != nil {
					utilityID := int(tariff.LseID)
					tariffID := int(tariff.ID)
					lead.UtilityID = &utilityID
					lead.TariffID = &tariffID

					if err := s.leadRepo.Update(ctx, &lead); err != nil {
						log.Printf("Warning: Failed to update lead with Genability data: %v", err)
					}
				}
			}
		}

	return &CreateLeadResponse{
		Success: true,
		LeadID:  lead.ID,
		HouseID: int(houseID),
	}, nil
}
