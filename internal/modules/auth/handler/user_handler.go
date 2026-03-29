package handler

import (
	"net/http"

	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/service"
	"github.com/NhatHaoDev3324/zizone-be/pkg/response"

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

	response.SuccessNoData(ctx, "Email verified successfully.")
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
		response.Fail(ctx, http.StatusInternalServerError, "Could not authenticate with Google: "+err.Error())
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

func (h *UserHandler) ForgotPassword(ctx *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.Email == "" {
		response.Fail(ctx, http.StatusBadRequest, "Email is required")
		return
	}

	if err := h.service.ForgotPassword(input.Email); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not send reset password email"+err.Error())
		return
	}

	response.SuccessNoData(ctx, "Reset password email sent successfully. Please check your email for the OTP code.")
}

func (h *UserHandler) VerifyOTPForgotPassword(ctx *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	token, err := h.service.VerifyOTPForgotPassword(input.Email, input.OTP)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithToken(ctx, "Email verified successfully.", token)
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	var input struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.NewPassword == "" {
		response.Fail(ctx, http.StatusBadRequest, "Password is required")
		return
	}

	if len(input.NewPassword) < 6 {
		response.Fail(ctx, http.StatusBadRequest, "Password must be at least 6 characters long")
		return
	}

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

	if err := h.service.ResetPassword(userID, input.NewPassword); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not reset password: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "Password reset successfully.")
}
