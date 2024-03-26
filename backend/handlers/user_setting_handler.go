package handlers

import (
	"dam/apis"
	"dam/enums"
	"dam/models"
	"dam/repositories"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSettingHandler struct {
	UserSettingRepo repositories.UserSettingRepoInterface
	db              *gorm.DB
}

type UserSettingHandlerInterface interface {
	CreateUserSetting(c *gin.Context)
}

func NewUserSettingHandler(db *gorm.DB) UserSettingHandlerInterface {
	return &UserSettingHandler{
		UserSettingRepo: repositories.NewUserSettingRepo(db),
		db:              db,
	}
}

func (h *UserSettingHandler) CreateUserSetting(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ctx.Value(enums.UserIDCtxKey).(string)

	var createUserSettingReq apis.CreateUserSettingRequest
	if err := c.BindJSON(&createUserSettingReq); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	if err := createUserSettingReq.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InvalidRequestError,
		})
		return
	}

	userSettingID := uuid.New().String()
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if _, err := repositories.GetUserSettingsByUserID(ctx, tx, userID, true); err == nil {
			return fmt.Errorf("user setting already exists")
		}

		userSetting := &models.UserSetting{
			UserSettingID: userSettingID,
			UserID:        userID,
			StorageVendor: createUserSettingReq.StorageVendor,
			StorageCredentials: &models.StorageCredentials{
				AWSS3AccessKeyID:     createUserSettingReq.AWSS3AccessKey,
				AWSS3SecretAccessKey: createUserSettingReq.AWSS3SecretKey,
			},
			StorageInformations: &models.StorageInformations{
				AWSS3BucketName: createUserSettingReq.AWSS3BucketName,
				AWSS3Region:     createUserSettingReq.AWSS3Region,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := repositories.CreateUserSetting(ctx, tx, userSetting); err != nil {
			return err
		}

		if createUserSettingReq.StorageVendor == string(enums.StorageAmazonS3) {
			sdkConfig, err := config.LoadDefaultConfig(
				ctx,
				config.WithRegion(createUserSettingReq.AWSS3Region),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
					createUserSettingReq.AWSS3AccessKey,
					createUserSettingReq.AWSS3SecretKey,
					""),
				),
			)
			if err != nil {
				return fmt.Errorf("load default config error: %w", err)
			}

			// Create S3 bucket
			s3Client := s3.NewFromConfig(sdkConfig)
			_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
				Bucket: &createUserSettingReq.AWSS3BucketName,
			})
			if err != nil {
				return fmt.Errorf("create bucket error: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusCreated, apis.CreateUserSettingResponse{
		UserSettingID: userSettingID,
	})
}
