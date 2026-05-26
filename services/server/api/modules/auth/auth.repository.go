package auth

import (
	queries "github.com/Shoyeb45/server/pkg/repository/gen-queries"
)

type AuthRepository struct {
	q *queries.Queries
}

func NewAuthRepository(q *queries.Queries) *AuthRepository {
    return &AuthRepository{q: q}
}
