package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev681999/tpo-platform/pkg/config"
	"github.com/dev681999/tpo-platform/pkg/db"
	"github.com/dev681999/tpo-platform/pkg/ent"
	"github.com/dev681999/tpo-platform/pkg/ent/user"
	"github.com/dev681999/tpo-platform/pkg/transport"
	userserver "github.com/dev681999/tpo-platform/pkg/user/server"
	usertransport "github.com/dev681999/tpo-platform/pkg/user/transport"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/run"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func checkErr(err error) {
	if err != nil {
		log.Err(err).Msg("exit")
		os.Exit(-1)
	}
}

type userService struct {
	db     db.DB
	logger zerolog.Logger
}

func (userService *userService) Get(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	u, err := userService.db.Client.User.Get(ctx, id)
	if err != nil {
		userService.logger.Debug().Err(err).Msg("")
		return nil, err
	}

	return u, nil
}

func main() {
	logger := zerolog.New(os.Stdout)
	logger = logger.With().Timestamp().Caller().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger = logger.With().Str("cmp", "global").Logger()

	log.Logger = logger

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	args := os.Args
	if len(args) > 1 {
		if args[1] == "init" {
			err := config.GenerateDefault()
			checkErr(err)
		}

		return
	}

	cfg, err := config.ReadConfig()
	checkErr(err)
	logger.Debug().Str("cfg", cfg.String()).Msg("")
	runServer(cfg, logger)
}

func runServer(cfg config.Config, logger zerolog.Logger) {
	db, err := db.New(cfg)
	checkErr(err)
	defer func() {
		logger.Info().Str("msg", "closing db").Msg("")
		if err := db.Close(); err != nil {
			logger.Fatal().Err(err).Msg("err while closing db")
		}
	}()

	e := transport.NewServer(logger)

	ug := e.Group("/user")

	us := userserver.New(db, logger.With().Str("cmp", "user").Str("layer", "server").Logger())

	usertransport.NewHandler(us, logger.With().Str("cmp", "user").Str("layer", "transport").Logger(), ug)

	us.GetAll(context.Background())

	var g run.Group
	{
		g.Add(func() error {
			logger.Info().Str("msg", "serving http").Msg("server")
			// return e.Start(":9000")
			return e.Start(cfg.ServerAddr)
		}, func(error) {
			logger.Info().Str("msg", "stopping http").Msg("server")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}
		})
	}
	{
		// set-up our signal handler
		var (
			cancelInterrupt = make(chan struct{})
			c               = make(chan os.Signal, 2)
		)
		defer close(c)

		g.Add(func() error {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Err(g.Run()).Msg("exit")
}

func getUsers(ctx context.Context, client *ent.Client, logger zerolog.Logger) error {
	users, err := client.User.Query().Order(ent.Desc(user.FieldCreateTime)).All(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		logger.Info().Str("user", u.String()).Msg("")
	}

	return nil
}

func createUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetName("a8m").
		SetPassword("123").
		SetEmail(fmt.Sprintf("%s@t.com", uuid.New().String())).
		SetRole(user.RoleCompany).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Printf("user was created: %+v", u)
	return u, nil
}
