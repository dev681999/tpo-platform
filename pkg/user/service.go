package user

import (
	"context"

	"github.com/dev681999/tpo-platform/pkg/ent"
	"github.com/dev681999/tpo-platform/pkg/ent/user"
	"github.com/google/uuid"
)

// Service is user service definition
type Service interface {
	Create(ctx context.Context, email, password, name string, role user.Role) (*ent.User, error)
	Get(ctx context.Context, id uuid.UUID) (*ent.User, error)
	GetAll(ctx context.Context) ([]*ent.User, error)
}
