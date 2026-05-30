package auth

import (
	"context"
	"errors"

	"github.com/Shoyeb45/server/pkg/apierr"
	queries "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/Shoyeb45/server/pkg/shared"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	q *queries.Queries
}

func NewAuthRepository(q *queries.Queries) *AuthRepository {
	return &AuthRepository{q: q}
}

func (r *AuthRepository) GetUserByGithubId(ctx context.Context, githubId int32) (*queries.User, error) {
	return r.q.GetUserByGithubId(ctx, githubId)
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID int32) (*queries.User, error) {
	user, err := r.q.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierr.NewNotFound("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, githubUser shared.GithubUser) (*queries.User, error) {
	return r.q.CreateUser(ctx, queries.CreateUserParams{
		Name:           githubUser.Name,
		GithubUsername: githubUser.Login,
		Email:          githubUser.Email,
		AvatarUrl:      githubUser.AvatarURL,
		GithubID:       int32(githubUser.ID),
	})
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID int32, rawToken string) error {
	_, err := r.q.CreateRefreshToken(ctx, queries.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: shared.HashToken(rawToken),
	})
	return err
}

func (r *AuthRepository) ValidateRefreshToken(ctx context.Context, rawToken string) (int32, error) {
	hash := shared.HashToken(rawToken)
	row, err := r.q.GetRefreshTokenByHash(ctx, hash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apierr.NewUnauthorized("refresh token revoked or invalid")
		}
		return 0, err
	}
	return row.UserID, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, rawToken string) error {
	return r.q.DeleteRefreshTokenByHash(ctx, shared.HashToken(rawToken))
}


func (r *AuthRepository) FindOrCreateUser(ctx context.Context, githubUser shared.GithubUser) (*queries.User, error) {
	user, err := r.GetUserByGithubId(ctx, int32(githubUser.ID))
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return r.CreateUser(ctx, githubUser)
}