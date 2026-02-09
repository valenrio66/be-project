package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/valenrio66/be-project/pkg/utils"

	"github.com/valenrio66/be-project/internal/db"
	"github.com/valenrio66/be-project/internal/dto"
)

var (
	ErrCampaignNotFound = errors.New("campaign not found")
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
		ID:          campaign.ID.String(),
		UserID:      campaign.UserID.String(),
		Title:       campaign.Title,
		Description: utils.PtrToString(campaign.Description),
		Status:      campaign.Status,
		Budget:      campaign.Budget,
		CreatedAt:   campaign.CreatedAt.Time,
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
			ID:          c.ID.String(),
			UserID:      c.UserID.String(),
			Title:       c.Title,
			Description: utils.PtrToString(c.Description),
			Status:      c.Status,
			StartDate:   c.StartDate.Time,
			EndDate:     c.EndDate.Time,
			Budget:      c.Budget,
			CreatedAt:   c.CreatedAt.Time,
		})
	}

	return responses, nil
}

func (s *CampaignService) GetCampaign(ctx context.Context, userID uuid.UUID, campaignID uuid.UUID) (*dto.CampaignResponse, error) {
	campaign, err := s.queries.GetCampaign(ctx, db.GetCampaignParams{
		ID:     campaignID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCampaignNotFound
		}
		return nil, err
	}

	var startDate, endDate time.Time
	if campaign.StartDate.Valid {
		startDate = campaign.StartDate.Time
	}
	if campaign.EndDate.Valid {
		endDate = campaign.EndDate.Time
	}

	return &dto.CampaignResponse{
		ID:          campaign.ID.String(),
		UserID:      campaign.UserID.String(),
		Title:       campaign.Title,
		Description: utils.PtrToString(campaign.Description),
		Status:      campaign.Status,
		Budget:      campaign.Budget,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   campaign.CreatedAt.Time,
	}, nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, userID uuid.UUID, campaignID uuid.UUID, req dto.UpdateCampaignRequest) (*dto.CampaignResponse, error) {
	var budget pgtype.Numeric
	if req.Budget != nil {
		if err := budget.Scan(*req.Budget); err != nil {
			return nil, err
		}
	}

	arg := db.UpdateCampaignParams{
		ID:          campaignID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		StartDate:   utils.ToPgTimestamp(req.StartDate),
		EndDate:     utils.ToPgTimestamp(req.EndDate),
		Budget:      budget,
	}

	campaign, err := s.queries.UpdateCampaign(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCampaignNotFound
		}
		return nil, err
	}

	var startDate, endDate time.Time
	if campaign.StartDate.Valid {
		startDate = campaign.StartDate.Time
	}
	if campaign.EndDate.Valid {
		endDate = campaign.EndDate.Time
	}

	return &dto.CampaignResponse{
		ID:          campaign.ID.String(),
		UserID:      campaign.UserID.String(),
		Title:       campaign.Title,
		Description: utils.PtrToString(campaign.Description),
		Status:      campaign.Status,
		Budget:      campaign.Budget,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   campaign.CreatedAt.Time,
	}, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, userID uuid.UUID, campaignID uuid.UUID) error {
	err := s.queries.DeleteCampaign(ctx, db.DeleteCampaignParams{
		ID:     campaignID,
		UserID: userID,
	})

	return err
}
