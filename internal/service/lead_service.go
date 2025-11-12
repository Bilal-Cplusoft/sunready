package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Bilal-Cplusoft/sunready/internal/client"
	"github.com/Bilal-Cplusoft/sunready/internal/models"
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
)

type CreateLeadResponse struct {
	Success bool `json:"success"`
	LeadID  int  `json:"lead_id"`
	HouseID int  `json:"house_id"`
}

type LeadService struct {
	leadRepo          *repo.LeadRepo
	projectRepo       *repo.ProjectRepo
	userRepo          *repo.UserRepo
	houseRepo         *repo.HouseRepo
	genabilityClient  *client.Agent
	lightFusionClient *client.LightFusionClient
}

type CreateLead struct {
	ProjectID         int     `json:"project_id" example:"1"`
	UserId            *int    `json:"user_id,omitempty" example:"1"`
	Latitude          float64 `json:"latitude" example:"37.7749"`
	Longitude         float64 `json:"longitude" example:"-122.4194"`
	SystemSize        float64 `json:"system_size" example:"10.5"`
	PanelCount        int     `json:"panel_count" example:"30"`
	HardwareType      *string `json:"hardware_type,omitempty"`
	KwhUsage          float64 `json:"kwh_usage" example:"12000"`
	Consumption       []int   `json:"consumption,omitempty"`
	LseId             int     `json:"lse_id"`
	Period            string  `json:"period"`
	TargetSolarOffset int     `json:"target_solar_offset"`
	Mode              *string `json:"mode,omitempty"`
	Unit              string  `json:"unit"`
}

func NewLeadService(leadRepo *repo.LeadRepo, houseRepo *repo.HouseRepo, lightFusionClient *client.LightFusionClient, projectRepo *repo.ProjectRepo, userRepo *repo.UserRepo) *LeadService {
	var genClient *client.Agent

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Failed to initialize Genability client: %v", r)
		}
	}()

	genClient = client.NewAgent()

	return &LeadService{
		leadRepo:          leadRepo,
		houseRepo:         houseRepo,
		genabilityClient:  genClient,
		projectRepo:       projectRepo,
		userRepo:          userRepo,
		lightFusionClient: lightFusionClient,
	}
}

func (s *LeadService) CreateLead(ctx context.Context, req CreateLead, userID int, projectID int) (*CreateLeadResponse, error) {
	if _, err := s.projectRepo.ExistsByID(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if _, err := s.userRepo.ExistsByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	house := models.House{
		Lat:   req.Latitude,
		Lng:   req.Longitude,
		State: "pending",
	}

	lead := models.Lead{
		ProjectID:  projectID,
		UserID:     &userID,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		KwhUsage:   req.KwhUsage,
		SystemSize: req.SystemSize,
		PanelCount: req.PanelCount,
		Consumption: req.Consumption,
		TargetSolarOffset: req.TargetSolarOffset,
	}
	User, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	addressDetails := client.AddressDetails{
		Street:     User.Street,
		City:       User.City,
		State:      User.State,
		PostalCode: User.PostalCode,
		Country:    User.Country,
	}
	hardware := client.HardwareDetails{
		PanelID: 156,
		InverterID: 324,
	}
	owner := client.HomeownerDetails{
		FirstName: User.FirstName,
		LastName:  User.LastName,
		Email:     User.Email,
		Phone:     User.PhoneNumber,
	}

	reqClient := client.Create3DProjectRequest{
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		Address:           addressDetails,
		TargetSolarOffset: req.TargetSolarOffset,
		Consumption:       req.Consumption,
        Unit: "kwh",
        Period: "year",
        Hardware: hardware,
        Homeowner: owner,
	}

	externalID, err := s.lightFusionClient.Create3DProject(ctx, reqClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create 3D project: %w", err)
	}
	if externalID == nil || externalID.LeadID == 0 {
		return nil, fmt.Errorf("failed to create 3D project: invalid response (nil or missing LeadID)")
	}

	lead.ExternalID = &externalID.LeadID

	if err := s.leadRepo.Create(ctx, &lead); err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}
	if err := s.houseRepo.Create(ctx, &house); err != nil {
		return nil, err
	}

	houseID := house.ID

	Lead, err := s.leadRepo.GetLeadWithUserByLeadID(ctx, lead.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead with user by user id: %w", err)
	}

	if s.genabilityClient != nil {
		accountInput := client.Account{
			Address: client.AccountAddress{
				String:    Lead.User.Street,
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

func (s *LeadService) GetMeshFiles(ctx context.Context, lead_external_id int) (*client.ProfilesFiles3DResponse, error) {
	resp, err := s.lightFusionClient.GetProjectFiles(ctx, lead_external_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get mesh files: %w", err)
	}
	return resp, nil
}
