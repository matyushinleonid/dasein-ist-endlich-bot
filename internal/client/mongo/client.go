package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Get(ctx context.Context, key int64, result interface{}) error
	Create(ctx context.Context, doc interface{}) (int64, error)
	Update(ctx context.Context, key int64, update interface{}) (int64, error)
	Delete(ctx context.Context, key int64) (int64, error)
}

type realClient struct {
	collection *mongo.Collection
}

func NewClient(cfg config.MongoDBConfig) (Client, error) {
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
	return &realClient{collection: coll}, nil
}

func (r *realClient) Get(ctx context.Context, key int64, result interface{}) error {
	filter := bson.M{"_id": key}
	err := r.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("document not found by key %d", key)
		}
		return err
	}
	return nil
}

func (r *realClient) Create(ctx context.Context, doc interface{}) (int64, error) {
	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return 0, err
	}
	id, ok := res.InsertedID.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected id type %T", res.InsertedID)
	}
	return id, nil
}

func (r *realClient) Update(ctx context.Context, key int64, update interface{}) (int64, error) {
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
	res, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": m})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *realClient) Delete(ctx context.Context, key int64) (int64, error) {
	filter := bson.M{"_id": key}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
