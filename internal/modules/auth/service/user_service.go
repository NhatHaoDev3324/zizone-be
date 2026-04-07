package service

import (
	"errors"
	"strings"

	"mime/multipart"

	"github.com/NhatHaoDev3324/zizone-be/config"
	"github.com/NhatHaoDev3324/zizone-be/constant"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/model"
	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/repository"
	"github.com/NhatHaoDev3324/zizone-be/tdo"
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
	//edit
	EditName(userID, fullName string) (string, error)
	EditPassword(userID, oldPassword, newPassword string) error
	//admin
	CreateAccount(fullName, email, role string) error
	GetAllUsers(page, limit int, search string) (tdo.Meta, []tdo.Profile, error)
	DeleteUser(id string) error
	EditAvatar(userID string, file *multipart.FileHeader) (string, error)
	GetDeletedUsers(page, limit int, search string) (tdo.Meta, []tdo.Profile, error)
	RestoreUser(id string) error
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

	token, err := utils.GenerateAccessToken(user.ID.String(), user.Role)
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

	token, err := utils.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	return s.repo.FindByIDNoCache(id)
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
	user, err := s.repo.FindByIDNoCache(userID)
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

func (s *userService) CreateAccount(fullName, email, role string) error {
	user, err := s.repo.FindByEmail(email)
	if err == nil {
		return errors.New("email already registered")
	}

	password, err := utils.SendPassword(email, fullName)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user = &model.User{
		ID:       uuid.New(),
		Email:    email,
		Password: hashedPassword,
		FullName: fullName,
		Avatar:   constant.NoAvatar,
		Provider: constant.ProviderAdminCreate,
		Role:     role,
		Active:   true,
	}

	if err := s.repo.Create(user); err != nil {
		return err
	}

	return err
}

func (s *userService) GetAllUsers(page, limit int, search string) (tdo.Meta, []tdo.Profile, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return tdo.Meta{}, nil, err
	}

	var result []tdo.Profile
	if search != "" {
		search = strings.ToLower(search)
		for _, profile := range users {
			fullName := strings.ToLower(profile.FullName)
			email := strings.ToLower(profile.Email)

			if strings.Contains(fullName, search) || strings.Contains(email, search) {
				result = append(result, profile)
			}
		}
	} else {
		result = users
	}

	total := len(result)
	totalPage := (total + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		return tdo.NewMetaResponse(total, totalPage, page, limit), []tdo.Profile{}, nil
	}
	if end > total {
		end = total
	}

	return tdo.NewMetaResponse(total, totalPage, page, limit), result[start:end], nil
}

func (s *userService) DeleteUser(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *userService) EditName(userID, fullName string) (string, error) {
	user, err := s.repo.FindByIDNoCache(userID)
	if err != nil {
		return "", err
	}

	user.FullName = fullName

	return user.FullName, s.repo.Update(user)
}

func (s *userService) EditPassword(userID, oldPassword, newPassword string) error {
	user, err := s.repo.FindByIDNoCache(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.repo.Update(user)
}

func (s *userService) EditAvatar(userID string, file *multipart.FileHeader) (string, error) {
	user, err := s.repo.FindByIDNoCache(userID)
	if err != nil {
		return "", err
	}

	url, err := utils.UploadImage(file)
	if err != nil {
		return "", err
	}

	user.Avatar = url

	return url, s.repo.Update(user)
}

func (s *userService) GetDeletedUsers(page, limit int, search string) (tdo.Meta, []tdo.Profile, error) {
	users, err := s.repo.FindAllDeleted()
	if err != nil {
		return tdo.Meta{}, nil, err
	}

	var result []tdo.Profile
	if search != "" {
		search = strings.ToLower(search)
		for _, profile := range users {
			fullName := strings.ToLower(profile.FullName)
			email := strings.ToLower(profile.Email)

			if strings.Contains(fullName, search) || strings.Contains(email, search) {
				result = append(result, profile)
			}
		}
	} else {
		result = users
	}

	total := len(result)
	totalPage := (total + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		return tdo.NewMetaResponse(total, totalPage, page, limit), []tdo.Profile{}, nil
	}
	if end > total {
		end = total
	}

	return tdo.NewMetaResponse(total, totalPage, page, limit), result[start:end], nil
}

func (s *userService) RestoreUser(id string) error {
	return s.repo.Restore(id)
}
