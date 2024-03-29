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
	FileVersionRepo repositories.FileVersionRepoInterface
}

type FileHandlerInterface interface {
	UploadFile(c *gin.Context)
	GetFile(c *gin.Context)
	UpdateFile(c *gin.Context)
	MoveFiles(c *gin.Context)
	ListFileVersions(c *gin.Context)
}

func NewFileHandler(db *gorm.DB) FileHandlerInterface {
	return &FileHandler{
		UserRepo:        repositories.NewUserRepo(db),
		DirectoryRepo:   repositories.NewDirectoryRepo(db),
		FileRepo:        repositories.NewFileRepo(db),
		FileVersionRepo: repositories.NewFileVersionRepo(db),
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

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}
	defer file.Close()

	directoryID := c.Param("directory_id")
	directory, err := h.DirectoryRepo.GetDirectoryByID(ctx, directoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "Directory not found",
			Code:    enums.DirectoryNotFoundError,
		})
		return
	}

	fileContentType := fileHeader.Header.Get("Content-Type")

	var fileM *models.File

	fileID := c.Query("file_id")
	switch fileID {
	case "":
		fileID = uuid.New().String()
		fileM = &models.File{
			FileID:      fileID,
			Name:        fileHeader.Filename,
			Size:        fileHeader.Size,
			Extension:   fileContentType,
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
	default:
		fileM, err = h.FileRepo.GetFileByID(ctx, fileID)
		if err != nil {
			c.JSON(http.StatusNotFound, apis.ErrorResponse{
				Message: "File not found",
				Code:    enums.FileNotFoundError,
			})
			return
		}
		fileID = fileM.FileID
	}

	fileVersion := &models.FileVersion{
		FileVersionID: uuid.New().String(),
		FileID:        fileID,
		Size:          fileHeader.Size,
		Extension:     fileContentType,
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := h.FileVersionRepo.CreateFileVersion(ctx, fileVersion); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	fileM.LatestFileVersionID = fileVersion.FileVersionID
	fileM.Size = fileHeader.Size
	fileM.UpdatedAt = time.Now()
	if err := h.FileRepo.UpdateFile(ctx, fileM); err != nil {
		c.JSON(http.StatusInternalServerError, apis.ErrorResponse{
			Message: err.Error(),
			Code:    enums.InternalError,
		})
		return
	}

	c.JSON(http.StatusCreated, apis.UploadFileResponse{
		FileVersionID: fileVersion.FileVersionID,
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

func (h *FileHandler) ListFileVersions(c *gin.Context) {
	ctx := c.Request.Context()

	fileID := c.Param("file_id")
	file, err := h.FileRepo.GetFileByID(ctx, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "File not found",
			Code:    enums.FileNotFoundError,
		})
		return
	}

	fileVersions, err := h.FileVersionRepo.ListFileVersions(ctx, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, apis.ErrorResponse{
			Message: "FileVersion not found",
			Code:    enums.FileVersionNotFoundError,
		})
		return
	}

	c.JSON(http.StatusOK, apis.ListFileVersions{
		File: apis.File{
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
		},
		FileVersions: func() []apis.FileVersion {
			fileVersionsAPI := make([]apis.FileVersion, 0, len(fileVersions))
			for _, fileVersion := range fileVersions {
				fileVersionsAPI = append(fileVersionsAPI, apis.FileVersion{
					FileVersionID: fileVersion.FileVersionID,
					FileID:        fileVersion.FileID,
					Size:          fileVersion.Size,
					Extension:     fileVersion.Extension,
					UserID:        fileVersion.UserID,
					CreatedAt:     fileVersion.CreatedAt,
					UpdatedAt:     fileVersion.UpdatedAt,
				})
			}
			return fileVersionsAPI
		}(),
	})
}
