package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/model"
	"github.com/NhatHaoDev3324/zizone-be/tdo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindAll() ([]tdo.Profile, error)
	FindByID(id string) (*model.User, error)
	FindByIDNoCache(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id string) error
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

	userKey := fmt.Sprintf("user:%s", user.ID.String())
	userData, _ := json.Marshal(user)
	r.redis.Set(ctx, userKey, userData, 30*time.Minute)

	r.redis.Del(ctx, "users:all")

	return nil
}

func (r *userRepository) FindAll() ([]tdo.Profile, error) {
	ctx := context.Background()
	var profiles []tdo.Profile

	cachedUsers, err := r.redis.Get(ctx, "users:all").Result()
	if err == nil {
		if json.Unmarshal([]byte(cachedUsers), &profiles) == nil {
			return profiles, nil
		}
	}

	err = r.db.Model(&model.User{}).Find(&profiles).Error
	if err == nil {
		usersData, _ := json.Marshal(profiles)
		r.redis.Set(ctx, "users:all", usersData, 30*time.Minute)
	}

	return profiles, err
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

func (r *userRepository) FindByIDNoCache(id string) (*model.User, error) {
	var user model.User

	err := r.db.Where("id = ?", id).First(&user).Error
	if err == nil {
		return &user, nil
	}

	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).Limit(1).Find(&user).Error
	if err != nil {
		return nil, err
	}
	if user.Email == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		return err
	}

	ctx := context.Background()
	userKey := fmt.Sprintf("user:%s", user.ID)
	r.redis.Del(ctx, userKey)
	r.redis.Del(ctx, "users:all")

	return nil
}

func (r *userRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	ctx := context.Background()
	userKey := fmt.Sprintf("user:%s", id)
	r.redis.Del(ctx, userKey)
	r.redis.Del(ctx, "users:all")

	return nil
}
