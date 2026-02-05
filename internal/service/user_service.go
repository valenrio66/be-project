package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/valenrio66/be-project/config"
	"github.com/valenrio66/be-project/internal/db"
	"github.com/valenrio66/be-project/internal/dto"
	"github.com/valenrio66/be-project/pkg/token"
)

var (
	ErrUserAlreadyExists  = errors.New("email already exists")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	queries    *db.Queries
	tokenMaker *token.JWTMaker
	config     config.Config
}

func NewUserService(q *db.Queries, tokenMaker *token.JWTMaker, cfg config.Config) *UserService {
	return &UserService{
		queries:    q,
		tokenMaker: tokenMaker,
		config:     cfg,
	}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	arg := db.CreateUserParams{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	user, err := s.queries.CreateUser(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID.String(), user.Email, user.Role, s.config.TokenDuration)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
