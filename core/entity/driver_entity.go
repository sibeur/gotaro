package entity

import (
	"time"

	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/common/driver"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Driver struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt    time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt    time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	Slug         string    `bson:"slug,omitempty" json:"slug,omitempty"`
	Name         string    `bson:"name,omitempty" json:"name,omitempty"`
	Type         uint32    `bson:"type,omitempty" json:"type,omitempty"`
	DriverConfig any       `bson:"driver_config,omitempty" json:"driver_config,omitempty"`
}

func (col *Driver) ToJSON() common.GotaroMap {

	return common.GotaroMap{
		"id":            col.ID,
		"created_at":    common.DateTimeNullableToString(&col.CreatedAt),
		"updated_at":    common.DateTimeNullableToString(&col.UpdatedAt),
		"deleted_at":    common.DateTimeNullableToString(&col.DeletedAt),
		"slug":          col.Slug,
		"name":          col.Name,
		"type":          col.Type,
		"driver_config": common.DToMap(col.DriverConfig.(primitive.D)),
	}
}

func (col *Driver) ToJSONSimple() common.GotaroMap {
	return common.GotaroMap{
		"id":   col.ID,
		"slug": col.Slug,
		"name": col.Name,
		"type": col.Type,
	}
}

func (col *Driver) GetDefaultFolder() string {
	folder := "/"
	switch col.Type {
	case uint32(driver.GCSDriverType):
		driverConfig := col.DriverConfig.(primitive.D)
		driverConfigMap := common.DToMap(driverConfig)
		if driverConfigMap["default_folder"] != nil {
			folder = driverConfigMap["default_folder"].(string)
		}
	}
	return folder
}

func (col *Driver) GetDriverConfig() map[string]any {
	return common.DToMap(col.DriverConfig.(primitive.D))

}

func (col *Driver) GetFilePathFromDriver(targetFilePath string) string {
	switch col.Type {
	case uint32(driver.GCSDriverType):
		driverConfig := col.DriverConfig.(primitive.D)
		driverConfigMap := common.DToMap(driverConfig)
		bucketName := ""
		if driverConfigMap["bucket_name"] != nil {
			bucketName = driverConfigMap["bucket_name"].(string)
		}
		if bucketName != "" {
			targetFilePath = "gs://" + bucketName + "/" + targetFilePath
		}
	}
	return targetFilePath
}

func (col Driver) GetCollName() string {
	return "drivers"
}
