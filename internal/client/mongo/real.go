package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/retry"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RealClient struct {
	collection *mongo.Collection
}

func NewRealClient(cfg config.MongoDBConfig) (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	if cfg.Username != "" && cfg.Password != "" {
		clientOpts.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}

	cli, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect failed: %w", err)
	}
	if err := cli.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	coll := cli.Database(cfg.Database).Collection(cfg.Collection)
	return &RealClient{collection: coll}, nil
}

func (r *RealClient) Get(ctx context.Context, key int64, result interface{}) error {
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		return r.collection.FindOne(ctx, bson.M{"_id": key}).Decode(result)
	})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *RealClient) Create(ctx context.Context, doc interface{}) (int64, error) {
	var res *mongo.InsertOneResult
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		var e error
		res, e = r.collection.InsertOne(ctx, doc)
		return e
	})
	if err != nil {
		return 0, err
	}
	id, ok := res.InsertedID.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected id type %T", res.InsertedID)
	}
	return id, nil
}

func (r *RealClient) Update(ctx context.Context, key int64, update interface{}) (int64, error) {
	raw, err := bson.Marshal(update)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal update: %w", err)
	}
	var m bson.M
	if err := bson.Unmarshal(raw, &m); err != nil {
		return 0, fmt.Errorf("failed to unmarshal update to map: %w", err)
	}
	delete(m, "_id")

	filter := bson.M{"_id": key}
	var res *mongo.UpdateResult
	err = retry.Do(ctx, retry.DefaultConfig(), func() error {
		var e error
		res, e = r.collection.UpdateOne(ctx, filter, bson.M{"$set": m})
		return e
	})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *RealClient) Delete(ctx context.Context, key int64) (int64, error) {
	filter := bson.M{"_id": key}
	var res *mongo.DeleteResult
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		var e error
		res, e = r.collection.DeleteOne(ctx, filter)
		return e
	})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *RealClient) FindAll(ctx context.Context) (Cursor, error) {
	var cursor *mongo.Cursor
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		var e error
		cursor, e = r.collection.Find(ctx, bson.M{}, options.Find().SetBatchSize(100))
		return e
	})
	if err != nil {
		return nil, fmt.Errorf("mongo FindAll failed: %w", err)
	}
	return &realCursor{cursor: cursor}, nil
}

type realCursor struct {
	cursor *mongo.Cursor
}

func (c *realCursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c *realCursor) Decode(v interface{}) error {
	return c.cursor.Decode(v)
}

func (c *realCursor) Err() error {
	return c.cursor.Err()
}

func (c *realCursor) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}
