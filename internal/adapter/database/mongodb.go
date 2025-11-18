package database

import (
	"context"
	"errors"
	"example-otelgo-miniservice/internal/domain"
	"example-otelgo-miniservice/internal/utils"
	"fmt"
	"time"

	"example-otelgo-miniservice/internal/port"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
	opts   domain.MongoOptions
}

func NewMongoDB(opts domain.MongoOptions) (port.MongoDB, error) {
	if opts.URL == "" || opts.DBName == "" {
		return nil, errors.New("MongoOptions.URI and DBName are required")
	}
	if opts.ConnTimeoutSec <= 0 {
		opts.ConnTimeoutSec = 10
	}
	if opts.OpTimeoutSec <= 0 {
		opts.OpTimeoutSec = 5
	}
	if opts.TxnTimeoutSec <= 0 {
		opts.TxnTimeoutSec = 30
	}

	var client *mongo.Client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	err := utils.Retry(context.Background(), utils.RetryConfig{
		MaxAttempts: 5,
		Interval:    3 * time.Second,
		Backoff:     true,
	}, func(attempt int) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.ConnTimeoutSec)*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(options.Client().
			ApplyURI(opts.URL).
			SetServerAPIOptions(serverAPI).
			SetServerSelectionTimeout(time.Duration(opts.ConnTimeoutSec) * time.Second),
		)
		if err != nil {
			return fmt.Errorf("mongo connect: %w", err)
		}

		if pingErr := client.Ping(ctx, readpref.Primary()); pingErr != nil {
			_ = client.Disconnect(context.Background())
			return fmt.Errorf("mongo ping: %w", pingErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	db := client.Database(opts.DBName)

	return &MongoDB{
		client: client,
		db:     db,
		opts:   opts,
	}, nil
}

func (m *MongoDB) Client() *mongo.Client { return m.client }
func (m *MongoDB) DB() *mongo.Database   { return m.db }

func (m *MongoDB) OpContext(parent context.Context) (context.Context, context.CancelFunc) {
	to := time.Duration(m.opts.OpTimeoutSec) * time.Second
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, to)
}

func (m *MongoDB) TxnContext(parent context.Context) (context.Context, context.CancelFunc) {
	to := time.Duration(m.opts.TxnTimeoutSec) * time.Second
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, to)
}

func (m *MongoDB) RunTxn(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)

	txnCtx, cancel := m.TxnContext(ctx)
	defer cancel()

	return mongo.WithSession(txnCtx, sess, func(scCtx context.Context) error {
		if err := sess.StartTransaction(); err != nil {
			return err
		}

		if err := fn(scCtx); err != nil {
			_ = sess.AbortTransaction(context.Background())
			return err
		}
		return sess.CommitTransaction(context.Background())
	})
}

func (m *MongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
