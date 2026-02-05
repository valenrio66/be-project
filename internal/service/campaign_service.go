package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/valenrio66/be-project/pkg/utils"

	"github.com/valenrio66/be-project/internal/db"
	"github.com/valenrio66/be-project/internal/dto"
)

type CampaignService struct {
	queries *db.Queries
}

func NewCampaignService(queries *db.Queries) *CampaignService {
	return &CampaignService{
		queries: queries,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, userID uuid.UUID, req dto.CreateCampaignRequest) (*dto.CampaignResponse, error) {
	arg := db.CreateCampaignParams{
		UserID:      userID,
		Title:       req.Title,
		Description: utils.StringToPtr(req.Description),
		Status:      "draft",
		StartDate:   pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:     pgtype.Timestamptz{Time: req.EndDate, Valid: true},
		Budget:      req.Budget,
	}

	campaign, err := s.queries.CreateCampaign(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &dto.CampaignResponse{
		ID:        campaign.ID.String(),
		UserID:    campaign.UserID.String(),
		Title:     campaign.Title,
		Status:    campaign.Status,
		Budget:    campaign.Budget,
		CreatedAt: campaign.CreatedAt.Time,
	}, nil
}

func (s *CampaignService) ListCampaigns(ctx context.Context, userID uuid.UUID, page, limit int) ([]dto.CampaignResponse, error) {
	offset := (page - 1) * limit

	arg := db.ListCampaignsParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	campaigns, err := s.queries.ListCampaigns(ctx, arg)
	if err != nil {
		return nil, err
	}

	var responses []dto.CampaignResponse
	for _, c := range campaigns {
		responses = append(responses, dto.CampaignResponse{
			ID:        c.ID.String(),
			Title:     c.Title,
			Status:    c.Status,
			StartDate: c.StartDate.Time,
			EndDate:   c.EndDate.Time,
			Budget:    c.Budget,
			CreatedAt: c.CreatedAt.Time,
		})
	}

	return responses, nil
}
