package auth

import (
	"context"

	queries "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/Shoyeb45/server/pkg/shared"
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

func (r *AuthRepository) CreateUser(ctx context.Context, githubUser shared.GithubUser) (*queries.User, error) {
	return r.q.CreateUser(ctx, queries.CreateUserParams{
		Name:           githubUser.Name,
		GithubUsername: githubUser.Login,
		Email:          githubUser.Email,
		AvatarUrl:      githubUser.AvatarURL,
		GithubID:       int32(githubUser.ID),
	})
}
