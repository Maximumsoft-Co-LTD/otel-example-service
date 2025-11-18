package port

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDB interface {
	Client() *mongo.Client
	DB() *mongo.Database

	OpContext(parent context.Context) (context.Context, context.CancelFunc)
	TxnContext(parent context.Context) (context.Context, context.CancelFunc)

	RunTxn(ctx context.Context, fn func(ctx context.Context) error) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}
