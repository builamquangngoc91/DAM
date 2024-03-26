package handlers

import (
	"dam/apis"
	"dam/enums"
	"dam/models"
	"dam/repositories"

	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileHandler struct {
	UserRepo        repositories.UserRepoInterface
	UserSettingRepo repositories.UserSettingRepoInterface
	DirectoryRepo   repositories.DirectoryRepoInterface
	FileRepo        repositories.FileRepoInterface
}

type FileHandlerInterface interface {
	UploadFile(c *gin.Context)
	GetFile(c *gin.Context)
	UpdateFile(c *gin.Context)
	MoveFiles(c *gin.Context)
}

func NewFileHandler(db *gorm.DB) FileHandlerInterface {
	return &FileHandler{
		UserRepo:      repositories.NewUserRepo(db),
		DirectoryRepo: repositories.NewDirectoryRepo(db),
		FileRepo:      repositories.NewFileRepo(db),
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ctx.Value(enums.UserIDCtxKey).(string)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.MissingFileError,
		})
		return
	}

	directoryID := c.Param("directory_id")
	directory, err := h.DirectoryRepo.GetDirectoryByID(ctx, directoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}
	defer file.Close()

	fileM := &models.File{
		FileID:      uuid.New().String(),
		Name:        fileHeader.Filename,
		Size:        fileHeader.Size,
		Extension:   fileHeader.Header.Get("Content-Type"),
		FullPath:    directory.FullPath + "/" + fileHeader.Filename,
		UserID:      userID,
		DirectoryID: directoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.FileRepo.CreateFile(ctx, fileM); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusCreated, apis.UploadFileResponse{
		FileID: fileM.FileID,
	})
}

func (h *FileHandler) GetFile(c *gin.Context) {
	ctx := c.Request.Context()

	fileID := c.Param("file_id")
	file, err := h.FileRepo.GetFileByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "File not found",
				Code:    enums.FileNotFoundError,
			})
			return
		}
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.FileNotFoundError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.File{
		FileID:      file.FileID,
		Name:        file.Name,
		Size:        file.Size,
		Extension:   file.Extension,
		UserID:      file.UserID,
		DirectoryID: file.DirectoryID,
		FullPath:    file.FullPath,
		Description: file.Description,
		Tags:        file.Tags,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
	})
}

func (h *FileHandler) UpdateFile(c *gin.Context) {
	ctx := c.Request.Context()

	var req apis.UpdateFileRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	fileID := c.Param("file_id")
	file, err := h.FileRepo.GetFileByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "File not found",
				Code:    enums.FileNotFoundError,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	file.Description = req.Description
	file.Tags = req.Tags
	file.UpdatedAt = time.Now()

	if err := h.FileRepo.UpdateFile(ctx, file); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *FileHandler) MoveFiles(c *gin.Context) {
	ctx := c.Request.Context()

	var req apis.MoveFilesRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.BindJSONError,
		})
		return
	}

	destinationDirectory, err := h.DirectoryRepo.GetDirectoryByID(ctx, req.DestinationDirectoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Destination directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	for _, fileID := range req.SourceFileIDs {
		file, err := h.FileRepo.GetFileByID(ctx, fileID)
		if err != nil {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "File not found",
				Code:    enums.FileNotFoundError,
			})
			return
		}

		textNeedReplaced := file.FullPath[0:strings.LastIndex(file.FullPath, "/")]
		file.FullPath = strings.ReplaceAll(file.FullPath, textNeedReplaced, destinationDirectory.FullPath)
		file.DirectoryID = destinationDirectory.DirectoryID
		file.UpdatedAt = time.Now()

		if err := h.FileRepo.UpdateFile(ctx, file); err != nil {
			c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
				Message: err.Error(),
				Code:    enums.InternalError,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
