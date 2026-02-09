package dto

import (
	"time"
)

type CreateCampaignRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
	Budget      float64   `json:"budget" binding:"required,gte=0"`
}

type CampaignResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Budget      float64   `json:"budget"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

type UpdateCampaignRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date" binding:"omitempty,gtfield=StartDate"`
	Budget      *float64   `json:"budget" binding:"omitempty,gte=0"`
}
