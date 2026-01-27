package handlers

import (
	"coderoot/lesson-03/examples/project/models"
	"coderoot/lesson-03/examples/project/services"
	"coderoot/lesson-03/examples/project/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
	jwtSecret   []byte
}

func NewUserHandler(userService *services.UserService, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		return
	}
	user, err := h.userService.CreateUser(req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		fmt.Println("Login c.ShouldBindJSON")
		return
	}
	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		fmt.Println("Login utils.ValidationError")
		return
	}

	token, err := utils.GenerateToken(h.jwtSecret, user.ID, user.Username)
	if err != nil {
		utils.HandleError(c, err)
		fmt.Println("Login utils.GenerateToken")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unathorized")
		fmt.Println("GetProfile  utils.Error")
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.HandleError(c, err)
		fmt.Println("GetProfile  h.userService.GetUserByID")
		return
	}
	utils.Success(c, models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized")
		fmt.Println("UpdateProfile  c.Get")
	}
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, parseValidationErrors(err))
		fmt.Println("UpdateProfile  c.ShouldBindJSON")
		return
	}
	user, err := h.userService.UpdateUser(userID.(uint), req)

	if err != nil {
		utils.HandleError(c, err)
		fmt.Println("UpdateProfile  h.userService.UpdateUser")
	}

	utils.Success(c, models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func parseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	// 简化处理，实际应该解析 binding 错误
	errors["general"] = err.Error()
	return errors
}
