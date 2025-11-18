package port

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CollectionNamer interface {
	CollName() string
}

type Repository[T CollectionNamer] interface {
	// Insert เอกสารใหม่
	Insert(
		ctx context.Context,
		doc *T,
		opts ...options.Lister[options.InsertOneOptions],
	) (*T, error)

	// FindOne: คืนเอกสารเดียว + found
	FindOne(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.FindOneOptions],
	) (*T, bool, error)

	// FindByID: sugar method
	FindByID(
		ctx context.Context,
		id any,
		opts ...options.Lister[options.FindOneOptions],
	) (*T, bool, error)

	// Find: คืน slice ของ T
	Find(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.FindOptions],
	) ([]T, error)

	UpdateOne(
		ctx context.Context,
		filter any,
		update any,
		opts ...options.Lister[options.UpdateOneOptions],
	) (*mongo.UpdateResult, error)

	UpdateMany(
		ctx context.Context,
		filter any,
		update any,
		opts ...options.Lister[options.UpdateManyOptions],
	) (*mongo.UpdateResult, error)

	// UpdateOneReturning: อัปเดตแล้วคืนเอกสารหลังอัปเดต
	UpdateOneReturning(
		ctx context.Context,
		filter any,
		update any,
		opts ...options.Lister[options.FindOneAndUpdateOptions],
	) (*T, bool, error)

	// DeleteOne: ลบเอกสาร
	DeleteOne(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.DeleteOneOptions],
	) (int64, error)

	// Count: นับเอกสาร
	Count(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.CountOptions],
	) (int64, error)

	// Coll: คืน *mongo.Collection low-level
	Coll() *mongo.Collection
}
