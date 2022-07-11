package database

import (
	"context"
	"time"

	"github.com/anicoll/unicom/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type Postgres struct {
	pool pool
}

type pool interface {
	Ping(ctx context.Context) error
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
}

func New(pool pool) *Postgres {
	return &Postgres{
		pool: pool,
	}
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) CreateCommunication(ctx context.Context, comm *model.Communication) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO communications (id, domain, "type")
		 VALUES ($1, $2, $3)`, comm.ID, comm.Domain, comm.Type)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	for _, channel := range comm.ResponseChannels {
		_, err := tx.Exec(ctx,
			`INSERT INTO response_channels (id, communication_id, "type", "url")
			 VALUES ($1, $2, $3, $4)`, channel.ID, comm.ID, channel.Type, channel.Url)
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *Postgres) SetCommunicationStatus(ctx context.Context, workflowId string, status model.Status, externalId *string) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE communications
		 SET "status" = $2, sent_at = $3, external_id = $4
		 WHERE id = $1`, workflowId, status, time.Now(), externalId)
	return err
}

func (p *Postgres) CreateResponseChannel(ctx context.Context, channel model.ResponseChannel) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO response_channels (id, communication_id, "type", "url")
		 VALUES ($1, $2, $3, $4)`, uuid.NewString(), channel.CommunicationID, channel.Type, channel.Url)
	return err
}

func (p *Postgres) SetResponseChannelStatus(ctx context.Context, id, externalId string, status model.Status) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE response_channels
		 SET "status" = $2, external_id = $3, sent_at = $4
		 WHERE id = $1`, id, status, externalId, time.Now())
	return err
}
