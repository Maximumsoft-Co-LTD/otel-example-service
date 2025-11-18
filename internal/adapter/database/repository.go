package database

import (
	"context"
	"errors"

	"example-otelgo-miniservice/internal/port"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository[T port.CollectionNamer] struct {
	db      *mongo.Database
	mongodb port.MongoDB
	coll    *mongo.Collection
}

func NewRepository[T port.CollectionNamer](mongodb port.MongoDB) (port.Repository[T], error) {
	var zero T
	collName := zero.CollName()
	return &Repository[T]{
		db:      mongodb.DB(),
		mongodb: mongodb,
		coll:    mongodb.DB().Collection(collName),
	}, nil
}

func (r *Repository[T]) Coll() *mongo.Collection { return r.coll }

func (r *Repository[T]) Insert(
	ctx context.Context,
	doc *T,
	opts ...options.Lister[options.InsertOneOptions],
) (*T, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	if _, err := r.coll.InsertOne(ctx, doc, opts...); err != nil {
		return nil, err
	}
	return doc, nil
}

func (r *Repository[T]) FindOne(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.FindOneOptions],
) (*T, bool, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	var out T
	err := r.db.Collection(r.coll.Name()).FindOne(ctx, filter, opts...).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &out, true, nil
}

func (r *Repository[T]) FindByID(
	ctx context.Context,
	id any,
	opts ...options.Lister[options.FindOneOptions],
) (*T, bool, error) {
	return r.FindOne(ctx, bson.M{"_id": id}, opts...)
}

func (r *Repository[T]) Find(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.FindOptions],
) ([]T, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	cur, err := r.db.Collection(r.coll.Name()).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		if err := cur.Close(ctx); err != nil {
			return
		}
	}(cur, ctx)

	var list []T
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository[T]) UpdateOne(
	ctx context.Context,
	filter any,
	update any,
	opts ...options.Lister[options.UpdateOneOptions],
) (*mongo.UpdateResult, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	res, err := r.db.Collection(r.coll.Name()).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository[T]) UpdateMany(
	ctx context.Context,
	filter any,
	update any,
	opts ...options.Lister[options.UpdateManyOptions],
) (*mongo.UpdateResult, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	res, err := r.db.Collection(r.coll.Name()).UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository[T]) UpdateOneReturning(
	ctx context.Context,
	filter any,
	update any,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) (*T, bool, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()

	base := options.FindOneAndUpdate().SetReturnDocument(options.After)
	all := append([]options.Lister[options.FindOneAndUpdateOptions]{base}, opts...)

	var out T
	err := r.db.Collection(r.coll.Name()).FindOneAndUpdate(ctx, filter, update, all...).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &out, true, nil
}

func (r *Repository[T]) DeleteOne(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.DeleteOneOptions],
) (int64, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	res, err := r.db.Collection(r.coll.Name()).DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *Repository[T]) Count(
	ctx context.Context,
	filter any,
	opts ...options.Lister[options.CountOptions],
) (int64, error) {
	ctx, cancel := r.mongodb.OpContext(ctx)
	defer cancel()
	return r.db.Collection(r.coll.Name()).CountDocuments(ctx, filter, opts...)
}
