package server

import (
	"context"

	password "github.com/dev681999/go-pass"
	"github.com/dev681999/tpo-platform/pkg/db"
	"github.com/dev681999/tpo-platform/pkg/ent"
	userent "github.com/dev681999/tpo-platform/pkg/ent/user"
	"github.com/dev681999/tpo-platform/pkg/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type userService struct {
	db     db.DB
	logger zerolog.Logger
	hash   password.Hash
}

// New returns a new user service server
func New(db db.DB, logger zerolog.Logger) user.Service {
	return &userService{
		db:     db,
		logger: logger,
	}
}

func (s *userService) Get(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	u, err := s.db.Client.User.Get(ctx, id)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return nil, err
	}

	return u, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*ent.User, error) {
	users, err := s.db.Client.User.Query().All(ctx)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return nil, err
	}

	return users, nil
}

func (s *userService) Create(ctx context.Context, email, password, name string, role userent.Role) (*ent.User, error) {
	pwd, err := s.hash.Generate(password)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return nil, err
	}

	u, err := s.db.Client.User.
		Create().
		SetEmail(email).
		SetPassword(pwd).
		SetName(name).
		SetRole(role).
		Save(ctx)

	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return nil, err
	}

	return u, nil
}
