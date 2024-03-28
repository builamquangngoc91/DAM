package apis

import (
	"dam/enums"

	"errors"
)

type CreateUserSettingRequest struct {
	StorageVendor   string `json:"storage_vendor"`
	AWSS3BucketName string `json:"aws_s3_bucket_name"`
	AWSS3Region     string `json:"aws_s3_region"`
	AWSS3AccessKey  string `json:"aws_s3_access_key"`
	AWSS3SecretKey  string `json:"aws_s3_secret_key"`
}

func (r *CreateUserSettingRequest) Validate() error {
	if r.StorageVendor == "" {
		return errors.New("storage_vendor is required")
	}

	switch r.StorageVendor {
	case string(enums.StorageAmazonS3):
		if r.AWSS3BucketName == "" {
			return errors.New("aws_s3_bucket_name is required")
		}
		if r.AWSS3Region == "" {
			return errors.New("aws_s3_region is required")
		}
		if r.AWSS3AccessKey == "" {
			return errors.New("aws_s3_access_key is required")
		}
		if r.AWSS3SecretKey == "" {
			return errors.New("aws_s3_secret_key is required")
		}
	default:
		return errors.New("storage_vendor is invalid")
	}

	return nil
}

type CreateUserSettingResponse struct {
	UserSettingID string `json:"user_setting_id"`
}
