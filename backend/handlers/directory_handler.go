package handlers

import (
	"dam/apis"
	"dam/enums"
	"dam/models"
	"dam/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DirectoryHandler struct {
	DirectoryRepo repositories.DirectoryRepoInterface
}

type DirectoryHandlerInterface interface {
	CreateDirectory(c *gin.Context)
	UpdateDirectory(c *gin.Context)
	GetDirectoryDetailsByID(c *gin.Context)
	GetDirectoryByID(c *gin.Context)
	ListFilesOrFoldersByDirectoryID(c *gin.Context)
}

func NewDirectoryHandler(db *gorm.DB) DirectoryHandlerInterface {
	return &DirectoryHandler{
		DirectoryRepo: repositories.NewDirectoryRepo(db),
	}
}

func (h *DirectoryHandler) CreateDirectory(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ctx.Value(enums.UserIDCtxKey).(string)

	var createDirReq apis.CreateDirectoryRequest
	if err := c.BindJSON(&createDirReq); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	parentDirectory, err := h.DirectoryRepo.GetDirectoryByID(ctx, createDirReq.ParentID)
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Parent directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	directionID := uuid.New().String()
	fullPath := parentDirectory.FullPath + "/" + parentDirectory.DirectoryID
	dir := &models.Directory{
		DirectoryID:       directionID,
		Name:              createDirReq.Name,
		UserID:            userID,
		FullPath:          fullPath,
		ParentDirectoryID: createDirReq.ParentID,
		Level:             parentDirectory.Level + 1,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := h.DirectoryRepo.CreateDirectory(ctx, dir); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusCreated, apis.UpdateDirectoryResponse{
		DirectoryID: dir.DirectoryID,
	})
}

func (h *DirectoryHandler) UpdateDirectory(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ctx.Value(enums.UserIDCtxKey).(string)

	var updateDirReq apis.UpdateDirectoryRequest
	if err := c.BindJSON(&updateDirReq); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	dir, err := h.DirectoryRepo.GetDirectoryByID(ctx, c.Param("directory_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	if dir.UserID != userID {
		c.JSON(http.StatusUnauthorized, apis.ErrorResponse{
			Message: "Insufficient permission",
			Code:    enums.InsufficientPermissionError,
		})
		return
	}

	dir.Name = updateDirReq.Name
	dir.UpdatedAt = time.Now()

	if err := h.DirectoryRepo.UpdateDirectory(ctx, dir); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.UpdateDirectoryResponse{
		DirectoryID: dir.DirectoryID,
	})
}

func (h *DirectoryHandler) GetDirectoryDetailsByID(c *gin.Context) {
	ctx := c.Request.Context()

	dir, err := h.DirectoryRepo.GetDirectoryByID(ctx, c.Param("directory_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.Directory{
		DirectoryID:       dir.DirectoryID,
		Name:              dir.Name,
		FullPath:          dir.FullPath,
		UserID:            dir.UserID,
		Level:             dir.Level,
		ParentDirectoryID: dir.ParentDirectoryID,
		CreatedAt:         dir.CreatedAt,
		UpdatedAt:         dir.UpdatedAt,
	})
}

func (h *DirectoryHandler) GetDirectoryByID(c *gin.Context) {
	ctx := c.Request.Context()

	dir, err := h.DirectoryRepo.GetDirectoryByID(ctx, c.Param("directory_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.Directory{
		DirectoryID:       dir.DirectoryID,
		Name:              dir.Name,
		FullPath:          dir.FullPath,
		UserID:            dir.UserID,
		Level:             dir.Level,
		ParentDirectoryID: dir.ParentDirectoryID,
		CreatedAt:         dir.CreatedAt,
		UpdatedAt:         dir.UpdatedAt,
	})
}

func (h *DirectoryHandler) ListFilesOrFoldersByDirectoryID(c *gin.Context) {
	ctx := c.Request.Context()

	dirID := c.Param("directory_id")

	orderByStr := c.Query("order_by")
	if orderByStr == "" {
		orderByStr = "created_at DESC"
	}
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: "Invalid limit",
			Code:    enums.InvalidRequestError,
		})
		return
	}
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: "Invalid offset",
			Code:    enums.InvalidRequestError,
		})
		return
	}

	if _, err := h.DirectoryRepo.GetDirectoryByID(ctx, dirID); err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	filesOrFolders, err := h.DirectoryRepo.ListFilesOrFoldersByDirectoryID(ctx, dirID, orderByStr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, filesOrFolders)
}
