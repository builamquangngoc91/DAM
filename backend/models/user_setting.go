package models

import "time"

type UserSetting struct {
	UserSettingID       string
	UserID              string
	StorageVendor       string
	StorageCredentials  *StorageCredentials  `gorm:"serializer:json"`
	StorageInformations *StorageInformations `gorm:"serializer:json"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type StorageInformations struct {
	AWSS3BucketName string
	AWSS3Region     string
}

type StorageCredentials struct {
	AWSS3AccessKeyID     string
	AWSS3SecretAccessKey string
}
