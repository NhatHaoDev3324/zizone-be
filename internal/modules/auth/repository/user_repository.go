package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindAll() ([]model.User, error)
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
}

type userRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{db, redis}
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}

	ctx := context.Background()

	userKey := fmt.Sprintf("user:%d", user.ID)
	userData, _ := json.Marshal(user)
	r.redis.Set(ctx, userKey, userData, 30*time.Minute)

	r.redis.Del(ctx, "users:all")

	return nil
}

func (r *userRepository) FindAll() ([]model.User, error) {
	ctx := context.Background()
	var users []model.User

	cachedUsers, err := r.redis.Get(ctx, "users:all").Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedUsers), &users) == nil {
			return users, nil
		}
	}

	err = r.db.Find(&users).Error
	if err == nil {
		usersData, _ := json.Marshal(users)
		r.redis.Set(ctx, "users:all", usersData, 30*time.Minute)
	}

	return users, err
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	ctx := context.Background()
	var user model.User
	userKey := fmt.Sprintf("user:%s", id)

	cachedUser, err := r.redis.Get(ctx, userKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			return &user, nil
		}
	}

	err = r.db.Where("id = ?", id).First(&user).Error
	if err == nil {
		userData, _ := json.Marshal(user)
		r.redis.Set(ctx, userKey, userData, 30*time.Minute)
	}

	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		return err
	}

	ctx := context.Background()
	userKey := fmt.Sprintf("user:%d", user.ID)
	userData, _ := json.Marshal(user)
	r.redis.Set(ctx, userKey, userData, 30*time.Minute)
	r.redis.Del(ctx, "users:all")

	return nil
}
