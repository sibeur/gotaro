package driver

import (
	"errors"

	"github.com/sibeur/gotaro/core/common"
)

type DriverClientUseCase interface {
	GetName() string
	GetType() StorageDriverType
	GetTypeString() string
	GetDriver() any
	UploadFile(filePath string, targetFilePath string) (string, error)
	GetSignedUrl(filePath string) (string, error)
	IsStorageAssetPublic() (bool, error)
	IsStorageBucketExist() (bool, error)
	ValidateDriver() error
	Close()
}

type DriverClient struct {
	driverName     string
	driverType     StorageDriverType
	isDriverPublic bool
	driverConfig   any
	driver         any
}

func NewDriverClient(driverName string, driverType StorageDriverType, driverConfig map[string]any) (*DriverClient, error) {
	var driver any
	var isDriverPublic bool
	switch driverType {
	case GCSDriverType:
		gcsDriver, err := NewGCPDriverClient(&GCSDriverConfig{
			ProjectID:      driverConfig["project_id"].(string),
			BucketName:     driverConfig["bucket_name"].(string),
			DefaultFolder:  driverConfig["default_folder"].(string),
			ServiceAccount: driverConfig["service_account"].(string),
		})
		if err != nil {
			return nil, err
		}
		isDriverPublic, _ = gcsDriver.IsStorageAssetPublic()
		driver = gcsDriver
	}
	return &DriverClient{
		driverName:     driverName,
		driverType:     driverType,
		isDriverPublic: isDriverPublic,
		driverConfig:   driverConfig,
		driver:         driver,
	}, nil
}

func (dc *DriverClient) GetName() string {
	return dc.driverName
}

func (dc *DriverClient) GetType() StorageDriverType {
	return dc.driverType
}

func (dc *DriverClient) GetTypeString() string {
	switch dc.driverType {
	case GCSDriverType:
		return "gcs"
	}
	return ""
}

func (dc *DriverClient) GetDriver() any {
	return dc.driver
}

func (dc *DriverClient) UploadFile(tempFilePath string, targetFilePath string) (string, error) {
	switch dc.driverType {
	case GCSDriverType:
		return dc.driver.(GCPDriverClientUseCase).UploadFile(tempFilePath, targetFilePath)
	}

	return "", nil
}

func (dc *DriverClient) GetSignedUrl(filePath string) (string, error) {
	switch dc.driverType {
	case GCSDriverType:
		return dc.driver.(GCPDriverClientUseCase).GetSignedUrl(filePath)
	}
	return "", nil
}

func (dc *DriverClient) IsStorageAssetPublic() (bool, error) {
	return dc.isDriverPublic, nil
}

func (dc *DriverClient) IsStorageBucketExist() (bool, error) {
	switch dc.driverType {
	case GCSDriverType:
		return dc.driver.(GCPDriverClientUseCase).IsStorageBucketExist()
	}
	return false, nil
}

func (dc *DriverClient) ValidateDriver() error {
	switch dc.driverType {
	case GCSDriverType:
		gcpDriver := dc.driver.(GCPDriverClientUseCase)
		if dc.driver == nil {
			return errors.New(common.ErrDriverNotInitiate)
		}

		errChan := make(chan error, 2)
		go func(gcpDriver GCPDriverClientUseCase, errChan chan error) {
			_, err := gcpDriver.IsHasStorageAdminPrivilage()

			if err != nil {
				if err.Error() == common.ErrBucketNotExistMsg {
					err = errors.New(common.ErrBucketNotExistMsg)
				}
				err = errors.New(common.ErrNotHaveStorageAdminPrivilageMsg)
			}

			if err != nil {
				errChan <- err
				return
			}
			errChan <- nil
		}(gcpDriver, errChan)
		go func(gcpDriver GCPDriverClientUseCase, errChan chan error) {
			_, err := gcpDriver.IsStorageBucketExist()
			if err != nil {
				errChan <- errors.New(common.ErrBucketNotExistMsg)
				return
			}
			errChan <- nil
		}(gcpDriver, errChan)

		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

func (dc *DriverClient) Close() {
	switch dc.driverType {
	case GCSDriverType:
		dc.driver.(GCPDriverClientUseCase).Close()
	}
}
