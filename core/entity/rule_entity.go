package entity

import (
	"time"

	"github.com/sibeur/gotaro/core/common"
)

type Rule struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	Slug      string    `bson:"slug,omitempty" json:"slug,omitempty"`
	Name      string    `bson:"name,omitempty" json:"name,omitempty"`
	MaxSize   uint64    `bson:"max_size,omitempty" json:"max_size,omitempty"`
	Mimes     []string  `bson:"mimes,omitempty" json:"mimes,omitempty"`
	DriverID  string    `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
}

func (col *Rule) ToJSON() common.GotaroMap {
	return common.GotaroMap{
		"id":         col.ID,
		"created_at": common.DateTimeNullableToString(&col.CreatedAt),
		"updated_at": common.DateTimeNullableToString(&col.UpdatedAt),
		"deleted_at": common.DateTimeNullableToString(&col.DeletedAt),
		"slug":       col.Slug,
		"name":       col.Name,
		"max_size":   col.MaxSize,
		"mimes":      col.Mimes,
		"driver_id":  col.DriverID,
	}
}

func (col *Rule) ToJSONSimple() common.GotaroMap {
	return common.GotaroMap{
		"id":        col.ID,
		"slug":      col.Slug,
		"name":      col.Name,
		"max_size":  col.MaxSize,
		"mimes":     col.Mimes,
		"driver_id": col.DriverID,
	}
}

func (col Rule) GetCollName() string {
	return "rules"
}
