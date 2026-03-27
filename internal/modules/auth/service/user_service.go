package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/NhatHaoDev3324/goAuth/constant"
	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/model"
	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/repository"
	"github.com/NhatHaoDev3324/goAuth/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterByEmail(firstName, lastName, email, password string) error
	RegisterByGoogle(code string) (string, error)
	LoginByEmail(email, password string) (string, error)
	VerifyOTP(email, otp string) error
	GetAllUsers() ([]model.User, error)
	GetUserByID(id string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) RegisterByEmail(firstName, lastName, email, password string) error {
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
		user.FirstName = firstName
		user.LastName = lastName
		user.FullName = firstName + " " + lastName
		if err := s.repo.Update(user); err != nil {
			return err
		}
	} else {
		user = &model.User{
			ID:        uuid.New(),
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  hashedPassword,
			FullName:  firstName + " " + lastName,
			Avatar:    constant.NoAvatar,
			Provider:  constant.ProviderEmail,
			Role:      constant.RoleUser,
			Active:    false,
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
	googleUser, err := s.getGoogleAccount(code)
	if err != nil {
		return "", err
	}
	user, err := s.repo.FindByEmail(googleUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &model.User{
				ID:        uuid.New(),
				FirstName: googleUser.FamilyName,
				LastName:  googleUser.GivenName,
				Email:     googleUser.Email,
				FullName:  googleUser.Name,
				Avatar:    googleUser.Picture,
				Provider:  constant.ProviderGoogle,
				Role:      constant.RoleUser,
				Active:    true,
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

type GoogleUser struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

func (s *userService) getGoogleAccount(code string) (*GoogleUser, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, errors.New("google credentials not configured")
	}

	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code for token: %v", resp.Status)
	}

	var tokenRes struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %v, %s", resp.Status, string(body))
	}

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
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
		return "", errors.New("account is not verified. please check your email for OTP")
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

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	return s.repo.FindByID(id)
}
