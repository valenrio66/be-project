package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/valenrio66/be-project/pkg/token"
	"go.uber.org/zap"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

type AuthPayload struct {
	UserID uuid.UUID
	Email  string
	Role   string
}

func AuthMiddleware(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			zap.L().Warn("Auth failed: no header provided", zap.String("ip", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is not provided"})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			zap.L().Warn("Auth failed: invalid format", zap.String("header", authorizationHeader))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			zap.L().Warn("Auth failed: unsupported auth type",
				zap.String("received_type", authorizationType),
				zap.String("ip", c.ClientIP()),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unsupported authorization type"})
			return
		}

		accessToken := fields[1]

		claims, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			zap.L().Warn("Auth failed: invalid token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access token is invalid or expired"})
			return
		}

		c.Set(AuthorizationPayloadKey, claims)

		c.Next()
	}
}

func GetAuthPayload(c *gin.Context) (*AuthPayload, error) {
	payload, exists := c.Get(AuthorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization payload is missing")
	}

	claims, ok := payload.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid authorization payload type")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email claim is missing")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("role claim is missing")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user_id claim is missing")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	return &AuthPayload{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}
