package service

import (
	"errors"

	"github.com/NhatHaoDev3324/zizone-be/config"
	"github.com/NhatHaoDev3324/zizone-be/constant"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/model"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/repository"
	"github.com/NhatHaoDev3324/zizone-be/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterByEmail(fullName, email, password string) error
	RegisterByGoogle(code string) (string, error)
	LoginByEmail(email, password string) (string, error)
	VerifyOTP(email, otp string) error
	GetUserByID(id string) (*model.User, error)
	ForgotPassword(email string) error
	VerifyOTPForgotPassword(email, otp string) (string, error)
	ResetPassword(userID, newPassword string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) RegisterByEmail(fullName, email, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user, err := s.repo.FindByEmail(email)
	if err == nil {
		if user.Active {
			return errors.New("email already registered and verified")
		}
		user.Password = hashedPassword
		user.FullName = fullName
		if err := s.repo.Update(user); err != nil {
			return err
		}
	} else {
		user = &model.User{
			ID:       uuid.New(),
			Email:    email,
			Password: hashedPassword,
			FullName: fullName,
			Avatar:   constant.NoAvatar,
			Provider: constant.ProviderEmail,
			Role:     constant.RoleUser,
			Active:   false,
		}
		if err := s.repo.Create(user); err != nil {
			return err
		}
	}

	_, err = utils.SendOTP(email, user.FullName)
	return err
}

func (s *userService) VerifyOTP(email, otp string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	valid, status, err := utils.VerifyOTP(email, otp)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New(status)
	}

	user.Active = true
	return s.repo.Update(user)
}

func (s *userService) RegisterByGoogle(code string) (string, error) {
	googleAuth, err := config.GetGoogleAuth(code)
	if err != nil {
		return "", err
	}
	user, err := s.repo.FindByEmail(googleAuth.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &model.User{
				ID:       uuid.New(),
				Email:    googleAuth.Email,
				FullName: googleAuth.FamilyName + " " + googleAuth.GivenName,
				Avatar:   googleAuth.Picture,
				Provider: constant.ProviderGoogle,
				Role:     constant.RoleUser,
				Active:   true,
			}
			if err := s.repo.Create(user); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	token, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) LoginByEmail(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if !user.Active {
		return "", errors.New("account is not verified")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) ForgotPassword(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	_, err = utils.SendOTP(email, user.FullName)
	return err
}

func (s *userService) VerifyOTPForgotPassword(email, otp string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	valid, status, err := utils.VerifyOTP(email, otp)
	if err != nil {
		return "", err
	}

	if !valid {
		return "", errors.New(status)
	}

	token, err := utils.GenerateResetPasswordToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) ResetPassword(userID, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Active = true

	return s.repo.Update(user)
}
