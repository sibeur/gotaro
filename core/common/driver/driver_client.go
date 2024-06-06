package driver

import (
	"log"
	"os"
)

type DriverClientUseCase interface {
	GetName() string
	GetType() StorageDriverType
	GetDriver() any
	UploadFile(filePath string, targetFilePath string) (string, error)
	Close()
}

type DriverClient struct {
	driverName   string
	driverType   StorageDriverType
	driverConfig any
	driver       any
}

func NewDriverClient(driverName string, driverType StorageDriverType, driverConfig map[string]any) (*DriverClient, error) {
	var driver any
	var err error
	switch driverType {
	case GCSDriverType:
		driver, err = NewGCPDriverClient(&GCSDriverConfig{
			ProjectID:      driverConfig["project_id"].(string),
			BucketName:     driverConfig["bucket_name"].(string),
			DefaultFolder:  driverConfig["default_folder"].(string),
			ServiceAccount: driverConfig["service_account"].(string),
		})
		if err != nil {
			return nil, err
		}
	}
	return &DriverClient{
		driverName:   driverName,
		driverType:   driverType,
		driverConfig: driverConfig,
		driver:       driver,
	}, nil
}

func (dc *DriverClient) GetName() string {
	return dc.driverName
}

func (dc *DriverClient) GetType() StorageDriverType {
	return dc.driverType
}

func (dc *DriverClient) GetDriver() any {
	return dc.driver
}

func (dc *DriverClient) UploadFile(tempFilePath string, targetFilePath string) (string, error) {
	defer func() {
		err := os.Remove(tempFilePath)
		if err != nil {
			log.Printf("Failed to delete temp file %v", tempFilePath)
		}
	}()

	switch dc.driverType {
	case GCSDriverType:
		return dc.driver.(GCPDriverClientUseCase).UploadFile(tempFilePath, targetFilePath)
	}

	return "", nil
}

func (dc *DriverClient) Close() {
	switch dc.driverType {
	case GCSDriverType:
		dc.driver.(GCPDriverClientUseCase).Close()
	}
}
