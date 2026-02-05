package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/valenrio66/be-project/internal/dto"
	"github.com/valenrio66/be-project/internal/middleware"
	"github.com/valenrio66/be-project/internal/service"
)

type CampaignHandler struct {
	campaignService *service.CampaignService
}

func NewCampaignHandler(campaignService *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
	}
}

// Create Campaign
// @Summary      Create new campaign
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateCampaignRequest true "Campaign Payload"
// @Success      201  {object}  dto.APIResponse{data=dto.CampaignResponse}
// @Failure      400  {object}  dto.APIResponse
// @Failure      500  {object}  dto.APIResponse
// @Router       /campaigns [post]
func (h *CampaignHandler) Create(c *gin.Context) {
	var req dto.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("CreateCampaign failed: invalid json", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.APIResponse{Error: err.Error()})
		return
	}

	authPayload, err := middleware.GetAuthPayload(c)
	if err != nil {
		zap.L().Warn("CreateCampaign failed: unauthorized", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dto.APIResponse{Error: "Unauthorized"})
		return
	}

	res, err := h.campaignService.CreateCampaign(c.Request.Context(), authPayload.UserID, req)
	if err != nil {
		zap.L().Error("CreateCampaign failed: service error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.APIResponse{Error: "Failed to create campaign"})
		return
	}

	zap.L().Info("Campaign created", zap.String("id", res.ID), zap.String("user_id", authPayload.UserID.String()))
	c.JSON(http.StatusCreated, dto.APIResponse{
		Message: "Campaign created successfully",
		Data:    res,
	})
}

// List Campaigns
// @Summary      List my campaigns
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Limit per page" default(10)
// @Success      200  {object}  dto.APIResponse{data=[]dto.CampaignResponse}
// @Router       /campaigns [get]
func (h *CampaignHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	authPayload, err := middleware.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{Error: "Unauthorized"})
		return
	}

	res, err := h.campaignService.ListCampaigns(c.Request.Context(), authPayload.UserID, page, limit)
	if err != nil {
		zap.L().Error("ListCampaigns failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.APIResponse{Error: "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Campaigns retrieved",
		Data:    res,
	})
}
