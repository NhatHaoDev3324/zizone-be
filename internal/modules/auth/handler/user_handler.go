package handler

import (
	"net/http"
	"strconv"

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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
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

	if input.Name == "" {
		response.Fail(ctx, http.StatusBadRequest, "Name is required")
		return
	}

	if err := h.service.RegisterByEmail(input.Email, input.Password, input.Name); err != nil {
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

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not get users")
		return
	}

	response.SuccessWithData(ctx, "Users fetched successfully", users)
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "User not found")
		return
	}

	response.SuccessWithData(ctx, "User fetched successfully", user)
}
