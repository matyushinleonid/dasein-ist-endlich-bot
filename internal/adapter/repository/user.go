package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/client/mongo"
	"github.com/matyushinleonid/dasein-ist-endlich-bot/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	mongoClient mongo.Client
}

func NewUserRepository(client mongo.Client) *UserRepository {
	return &UserRepository{mongoClient: client}
}

func (r *UserRepository) Get(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	if err := r.mongoClient.Get(ctx, id, &u); err != nil {
		if errors.Is(err, mongo.ErrNotFound) {
			return nil, fmt.Errorf("%w: user %d", ErrUserNotFound, id)
		}
		return nil, fmt.Errorf("db: could not get user %d: %w", id, err)
	}
	return &u, nil
}

func (r *UserRepository) UserExists(ctx context.Context, id int64) (bool, error) {
	var u model.User
	err := r.mongoClient.Get(ctx, id, &u)
	if err != nil {
		if errors.Is(err, mongo.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("db: checking existence for user %d: %w", id, err)
	}
	return true, nil
}

func (r *UserRepository) Create(ctx context.Context, id int64) (int64, error) {
	u := model.NewUser(id)
	newID, err := r.mongoClient.Create(ctx, *u)
	if err != nil {
		return 0, fmt.Errorf("db: could not create user: %w", err)
	}
	return newID, nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	modified, err := r.mongoClient.Update(ctx, u.ID, *u)
	if err != nil {
		return fmt.Errorf("db: could not update user %d: %w", u.ID, err)
	}
	if modified == 0 {
		return fmt.Errorf("%w: user %d", ErrUserNotFound, u.ID)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	deleted, err := r.mongoClient.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("db: could not delete user %d: %w", id, err)
	}
	if deleted == 0 {
		return fmt.Errorf("%w: user %d", ErrUserNotFound, id)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	cursor, err := r.mongoClient.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("db: could not list users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []model.User
	for cursor.Next(ctx) {
		var u model.User
		if err := cursor.Decode(&u); err != nil {
			return nil, fmt.Errorf("db: decode error: %w", err)
		}
		users = append(users, u)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("db: cursor error: %w", err)
	}
	return users, nil
}
