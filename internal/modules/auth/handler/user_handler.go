package handler

import (
	"net/http"

	"github.com/NhatHaoDev3324/goAuth/internal/modules/auth/service"
	"github.com/NhatHaoDev3324/goAuth/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) RegisterByEmail(ctx *gin.Context) {
	var input struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.FirstName == "" {
		response.Fail(ctx, http.StatusBadRequest, "First name is required")
		return
	}

	if input.LastName == "" {
		response.Fail(ctx, http.StatusBadRequest, "Last name is required")
		return
	}

	if input.Email == "" {
		response.Fail(ctx, http.StatusBadRequest, "Email is required")
		return
	}

	if input.Password == "" {
		response.Fail(ctx, http.StatusBadRequest, "Password is required")
		return
	}

	if len(input.Password) < 6 {
		response.Fail(ctx, http.StatusBadRequest, "Password must be at least 6 characters long")
		return
	}

	if err := h.service.RegisterByEmail(input.FirstName, input.LastName, input.Email, input.Password); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not register user: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "Registration successful. Please check your email for the OTP code.")
}

func (h *UserHandler) VerifyOTP(ctx *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.service.VerifyOTP(input.Email, input.OTP); err != nil {
		response.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessNoData(ctx, "Email verified successfully. You can now log in.")
}

func (h *UserHandler) LoginByEmail(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	token, err := h.service.LoginByEmail(input.Email, input.Password)
	if err != nil {
		if err.Error() == "account is not verified. please check your email for OTP" {
			response.Fail(ctx, http.StatusForbidden, err.Error())
			return
		}
		response.Fail(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	response.SuccessWithToken(ctx, "Logged in successfully", token)
}

func (h *UserHandler) RegisterByGoogle(ctx *gin.Context) {
	var input struct {
		Code string `json:"code" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	token, err := h.service.RegisterByGoogle(input.Code)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not authenticate with Google")
		return
	}

	response.SuccessWithToken(ctx, "Logged in successfully", token)
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		response.Fail(ctx, http.StatusUnauthorized, "User not found")
		return
	}
	userID, ok := userId.(string)
	if !ok {
		response.Fail(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "User not found")
		return
	}

	response.SuccessWithData(ctx, "User fetched successfully", user)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not get users")
		return
	}

	response.SuccessWithData(ctx, "Users fetched successfully", users)
}
