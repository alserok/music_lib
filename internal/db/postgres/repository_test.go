package postgres

import (
	"context"
	"github.com/alserok/music_lib/internal/config"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

type RepositorySuite struct {
	suite.Suite

	repo      *repository
	conn      *sqlx.DB
	container *postgres.PostgresContainer
}

func (suite *RepositorySuite) SetupTest() {
	suite.conn, suite.container = newPostgresDB(&suite.Suite)
	suite.repo = NewRepository(suite.conn)
}

func (suite *RepositorySuite) TeardownTest() {
	suite.Require().NoError(suite.conn.Close())
	suite.Require().NoError(suite.container.Terminate(context.Background()))
}

//func (suite *RepositorySuite) TestGetSongs() {
//
//}

func newPostgresDB(s *suite.Suite) (*sqlx.DB, *postgres.PostgresContainer) {
	ctx := context.Background()
	cfg := config.Postgres{
		Name: "postgres",
		Port: "5432",
		User: "postgres",
		Pass: "postgres",
		Host: "localhost",
	}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(cfg.Name),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	s.Require().NoError(err)
	s.Require().NotNil(postgresContainer)
	s.Require().True(postgresContainer.IsRunning())

	port, err := postgresContainer.MappedPort(ctx, nat.Port(cfg.Port+"/tcp"))
	s.Require().NoError(err)
	cfg.Port = port.Port()

	conn := MustConnect(cfg.DSN())

	return conn, postgresContainer
}
