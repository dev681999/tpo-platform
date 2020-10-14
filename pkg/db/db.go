package db

import (
	"context"
	"fmt"

	"github.com/dev681999/tpo-platform/pkg/config"
	"github.com/dev681999/tpo-platform/pkg/ent"
	"github.com/rs/zerolog/log"
)

// DB hold database data
type DB struct {
	Client *ent.Client
}

// New returns a new DB connection
func New(cfg config.Config) (DB, error) {
	client, err := ent.Open("sqlite3", cfg.DBFile)
	if err != nil {
		log.Debug().Err(err).Msg("")
		return DB{}, err
	}

	err = client.Schema.Create(context.Background())
	if err != nil {
		verr := client.Close()
		if verr != nil {
			err = fmt.Errorf("%s : %w", err.Error(), verr)
		}
		log.Debug().Err(err).Msg("")
		return DB{}, err
	}

	return DB{
		Client: client,
	}, nil
}

// Close close client connection
func (db DB) Close() error {
	return db.Client.Close()
}
