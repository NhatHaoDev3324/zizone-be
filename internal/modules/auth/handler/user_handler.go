package handler

import (
	"net/http"

	"github.com/NhatHaoDev3324/zizone-be/internal/modules/auth/service"
	"github.com/NhatHaoDev3324/zizone-be/pkg/response"
	"github.com/NhatHaoDev3324/zizone-be/tdo"

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
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.FullName == "" {
		response.Fail(ctx, http.StatusBadRequest, "Full name is required")
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

	if err := h.service.RegisterByEmail(input.FullName, input.Email, input.Password); err != nil {
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

	response.SuccessDataInfo(ctx, "User fetched successfully", tdo.NewProfile(user.ID.String(), user.Email, user.FullName, user.Avatar, user.Role, user.Provider, user.CreatedAt.String(), ""))
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

func (h *UserHandler) CreateAccount(ctx *gin.Context) {
	var input struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.FullName == "" {
		response.Fail(ctx, http.StatusBadRequest, "Full name is required")
		return
	}

	if input.Email == "" {
		response.Fail(ctx, http.StatusBadRequest, "Email is required")
		return
	}

	if input.Role == "" {
		response.Fail(ctx, http.StatusBadRequest, "Role is required")
		return
	}

	if err := h.service.CreateAccount(input.FullName, input.Email, input.Role); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not register user: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "Registration successful. Please check your email for the password.")
}

func (h *UserHandler) GetListUser(ctx *gin.Context) {
	var params struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Search string `form:"search"`
	}

	if params.Page == 0 {
		params.Page = 1
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	meta, users, err := h.service.GetAllUsers(params.Page, params.Limit, params.Search)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not get list user: "+err.Error())
		return
	}

	response.SuccessWithMetaAndData(ctx, "Get list user successfully", meta, users)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "" {
		response.Fail(ctx, http.StatusBadRequest, "User ID is required")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not delete user: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "User deleted successfully")
}

func (h *UserHandler) EditName(ctx *gin.Context) {
	var input struct {
		FullName string `json:"full_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.FullName == "" {
		response.Fail(ctx, http.StatusBadRequest, "Full name is required")
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

	userName, err := h.service.EditName(userID, input.FullName)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not edit name: "+err.Error())
		return
	}

	response.SuccessWithData(ctx, "Name edited successfully", tdo.NewProfile("", "", userName, "", "", "", "", ""))
}

func (h *UserHandler) EditPassword(ctx *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	if input.OldPassword == "" {
		response.Fail(ctx, http.StatusBadRequest, "Old password is required")
		return
	}

	if input.NewPassword == "" {
		response.Fail(ctx, http.StatusBadRequest, "New password is required")
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

	if err := h.service.EditPassword(userID, input.OldPassword, input.NewPassword); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not edit password: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "Password edited successfully")
}

func (h *UserHandler) EditAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Avatar file is required")
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

	url, err := h.service.EditAvatar(userID, file)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not edit avatar: "+err.Error())
		return
	}

	response.SuccessWithData(ctx, "Avatar edited successfully", tdo.NewProfile("", "", "", url, "", "", "", ""))
}

func (h *UserHandler) GetDeletedUsers(ctx *gin.Context) {
	var params struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Search string `form:"search"`
	}

	if params.Page == 0 {
		params.Page = 1
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "Invalid request")
		return
	}

	meta, users, err := h.service.GetDeletedUsers(params.Page, params.Limit, params.Search)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not get list deleted user: "+err.Error())
		return
	}

	response.SuccessWithMetaAndData(ctx, "Get list deleted user successfully", meta, users)
}

func (h *UserHandler) RestoreUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.Fail(ctx, http.StatusBadRequest, "User ID is required")
		return
	}

	if err := h.service.RestoreUser(id); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "Could not restore user: "+err.Error())
		return
	}

	response.SuccessNoData(ctx, "User restored successfully")
}
