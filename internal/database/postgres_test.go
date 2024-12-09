package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/anicoll/unicom/internal/database"
	"github.com/anicoll/unicom/internal/model"
)

type PostgresUnitTestSuite struct {
	suite.Suite
	dbdsn             string
	ctx               context.Context
	conn              *pgxpool.Pool
	postgresContainer testcontainers.Container
	postgres          *database.Postgres
}

func TestPostgresUnitTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresUnitTestSuite))
}

func (s *PostgresUnitTestSuite) SetupSuite() {
	s.ctx = context.Background()
	waiter := wait.ForLog("ready to accept connections")
	waiter.Occurrence = 2
	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: waiter,
	}
	var err error
	s.postgresContainer, err = testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.NoError(err)

	containerEndpoint, err := s.postgresContainer.Endpoint(s.ctx, "")
	s.NoError(err)
	s.dbdsn = fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", containerEndpoint)

	s.conn, err = pgxpool.Connect(s.ctx, s.dbdsn)
	s.NoError(err)

	s.postgres = database.New(s.conn, zap.NewNop())
	migrations := database.NewMigrations(s.dbdsn, database.MigrateUp)
	err = migrations.Execute()
	s.NoError(err)
}

func (s *PostgresUnitTestSuite) TearDownSuite() {
	err := s.postgresContainer.Terminate(s.ctx)
	s.NoError(err)
}

func (s *PostgresUnitTestSuite) Test_Ping_Success() {
	ctx := context.Background()

	err := s.postgres.Ping(ctx)
	s.NoError(err)
}

func (s *PostgresUnitTestSuite) Test_CreateCommunication_Success() {
	ctx := context.Background()

	expectedCommRequest := model.Communication{
		Domain: "test-domain",
		Status: model.Pending,
		Type:   model.Email,
	}

	err := s.postgres.CreateCommunication(ctx, &expectedCommRequest)
	s.NoError(err)

	row := s.conn.QueryRow(ctx,
		`SELECT id, domain, external_id, "type", created_at, sent_at, status
		FROM communications
		LIMIT 1; 
	`)

	got := model.Communication{}
	err = row.Scan(
		&got.ID,
		&got.Domain,
		&got.ExternalId,
		&got.Type,
		&got.CreatedAt,
		&got.SentAt,
		&got.Status,
	)
	s.NoError(err)
	if !cmp.Equal(expectedCommRequest, got, cmpopts.IgnoreFields(model.Communication{}, "CreatedAt", "SentAt", "ID", "ResponseChannels")) {
		s.Fail("expected not equal", cmp.Diff(expectedCommRequest, got, cmpopts.IgnoreFields(model.Communication{}, "CreatedAt", "SentAt", "ID", "ResponseChannels")))
	}
}
