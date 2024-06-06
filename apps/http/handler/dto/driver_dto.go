package dto

import "encoding/json"

type NewDriverDTO struct {
	Slug         string `json:"slug" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Type         uint32 `json:"type" validate:"required"`
	DriverConfig any    `json:"driver_config" validate:"required"`
}

type EditDriverDTO struct {
	Name         string `json:"name" validate:"required"`
	Type         uint32 `json:"type" validate:"required"`
	DriverConfig any    `json:"driver_config" validate:"required"`
}

type GCSDriverConfigDTO struct {
	ProjectID      string `json:"project_id" validate:"required"`
	BucketName     string `json:"bucket_name" validate:"required"`
	DefaultFolder  string `json:"default_folder" validate:"required"`
	ServiceAccount any    `json:"service_account" validate:"required"`
}

func (gdc *GCSDriverConfigDTO) GetServiceAccountJSONBytes() ([]byte, error) {
	return json.Marshal(gdc.ServiceAccount)
}
