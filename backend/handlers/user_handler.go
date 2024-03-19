package handlers

import (
	"dam/apis"
	"dam/enums"
	"dam/models"
	"dam/repositories"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	UserRepo repositories.UserRepoInterface
	RdClient *redis.Client
}

type UserHandlerInterface interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetCurrentUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

func NewUserHandler(db *gorm.DB, rdClient *redis.Client) UserHandlerInterface {
	return &UserHandler{
		UserRepo: repositories.NewUserRepo(db),
		RdClient: rdClient,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var loginReq apis.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	user, err := h.UserRepo.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "user not found",
				Code:    enums.UserNotFoundError,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: "wrong password",
			Code:    enums.InvalidRequestError,
		})
		return
	}

	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.UserID
	claims["exp"] = time.Now().Add(time.Hour).Unix() // Token expires in 1 hour

	// Generate the token string
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return
	}

	if err := h.RdClient.Set(ctx, tokenString, user.UserID, time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.LoginResponse{
		Token:     tokenString,
		Type:      "Bearer",
		ExpiresIn: 3600,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, apis.ErrorResponse{
			Message: "Authorization header is required",
			Code:    enums.AuthorizationHeaderRequiredError,
		})
		return
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, apis.ErrorResponse{
			Message: "Invalid authorization header",
			Code:    enums.InvalidAuthorizationHeaderError,
		})
		return
	}

	token := authParts[1]

	if err := h.RdClient.Del(ctx, token).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var createUserReq apis.CreateUserRequest
	if err := c.BindJSON(&createUserReq); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	if err := createUserReq.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InvalidRequestError,
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.HashPasswordError,
		})
		return
	}

	user := models.User{
		UserID:       uuid.New().String(),
		Username:     createUserReq.Username,
		PasswordHash: string(passwordHash),
		Email:        createUserReq.Email,
		Name:         createUserReq.Name,
	}

	if err := h.UserRepo.CreateUser(ctx, &user); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusCreated, apis.CreateUserResponse{
		UserID: user.UserID,
	})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ctx.Value(enums.UserIDCtxKey).(string)
	user, err := h.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "user not found",
				Code:    enums.UserNotFoundError,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.GetUserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ctx.Value(enums.UserIDCtxKey).(string)
	var updateUserReq apis.UpdateUserRequest
	if err := c.BindJSON(&updateUserReq); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	user, err := h.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "user not found",
				Code:    enums.UserNotFoundError,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	user.Name = updateUserReq.Name

	if err := h.UserRepo.UpdateUser(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.UpdateUserResponse{
		UserID: user.UserID,
	})
}
