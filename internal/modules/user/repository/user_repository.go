package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"template/internal/modules/user/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
}

type userRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
	return &userRepository{db, rdb}
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Cache the new user
	userKey := fmt.Sprintf("user:%d", user.ID)
	userData, _ := json.Marshal(user)
	r.rdb.Set(ctx, userKey, userData, 10*time.Minute)

	// Invalidate the "all_users" list cache
	r.rdb.Del(ctx, "users:all")

	return nil
}

func (r *userRepository) FindAll() ([]model.User, error) {
	ctx := context.Background()
	var users []model.User

	// Try to get from Redis
	cachedUsers, err := r.rdb.Get(ctx, "users:all").Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedUsers), &users) == nil {
			return users, nil
		}
	}

	// If not in Redis, get from database
	err = r.db.Find(&users).Error
	if err == nil {
		// Cache the results
		usersData, _ := json.Marshal(users)
		r.rdb.Set(ctx, "users:all", usersData, 10*time.Minute)
	}

	return users, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	ctx := context.Background()
	var user model.User
	userKey := fmt.Sprintf("user:%d", id)

	// Try to get from Redis
	cachedUser, err := r.rdb.Get(ctx, userKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			return &user, nil
		}
	}

	// If not in Redis, get from database
	err = r.db.First(&user, id).Error
	if err == nil {
		// Cache the result
		userData, _ := json.Marshal(user)
		r.rdb.Set(ctx, userKey, userData, 10*time.Minute)
	}

	return &user, err
}
