package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-project/internal/dto"
	"github.com/valenrio66/be-project/internal/middleware"
	"github.com/valenrio66/be-project/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register User
// @Summary      Register new user
// @Description  Add User
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register Payload"
// @Success      201  {object}  dto.APIResponse
// @Failure      400  {object}  dto.APIResponse
// @Failure      409  {object}  dto.APIResponse
// @Failure      500  {object}  dto.APIResponse
// @Router       /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("Register failed: invalid json", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.APIResponse{Error: err.Error()})
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			zap.L().Warn("Register failed: email duplicate", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, dto.APIResponse{Error: "Email already exists"})
			return
		}
		zap.L().Error("Register failed: system error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.APIResponse{Error: "Internal server error"})
		return
	}

	zap.L().Info("User registered successfully", zap.String("email", res.Email))
	c.JSON(http.StatusCreated, dto.APIResponse{
		Message: "User registered successfully",
		Data:    res,
	})
}

// Login User
// @Summary      Login user
// @Description  Login using email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login Credentials"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  dto.APIResponse
// @Failure      401  {object}  dto.APIResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Warn("Login failed: invalid json", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.APIResponse{Error: err.Error()})
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			zap.L().Warn("Login failed: invalid credentials", zap.String("email", req.Email))
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Error: "Invalid email or password"})
			return
		}

		zap.L().Error("Login failed: system error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.APIResponse{Error: "Internal server error"})
		return
	}

	zap.L().Info("User logged in", zap.String("email", res.User.Email))
	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Login successful",
		Data:    res,
	})
}

// GetMe
// @Summary      Get My Profile
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.APIResponse
// @Router       /me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	authPayload, err := middleware.GetAuthPayload(c)
	if err != nil {
		zap.L().Error("GetMe failed: auth payload error", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dto.APIResponse{Error: "Unauthorized"})
		return
	}

	user, err := h.service.GetUserByEmail(c.Request.Context(), authPayload.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.APIResponse{Error: "User not found"})
			return
		}
		zap.L().Error("GetMe failed: db error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.APIResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "User profile retrieved",
		Data:    user,
	})
}
