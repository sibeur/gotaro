package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sibeur/gotaro/core/common"
)

type Media struct {
	ID                 string    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt          time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt          time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt          time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	RuleSlug           string    `bson:"rule_slug,omitempty" json:"rule_slug,omitempty"`
	DriverSlug         string    `bson:"driver_slug,omitempty" json:"driver_slug,omitempty"`
	FilePath           string    `bson:"file_path,omitempty" json:"file_path,omitempty"`
	FilePathFromDriver string    `bson:"file_path_from_driver,omitempty" json:"file_path_from_driver,omitempty"`
	FileOriginalName   string    `bson:"file_original_name,omitempty" json:"file_original_name,omitempty"`
	FileAliasName      string    `bson:"file_alias_name,omitempty" json:"file_alias_name,omitempty"`
	FileSize           uint64    `bson:"file_size,omitempty" json:"file_size,omitempty"`
	FileMime           string    `bson:"file_mime,omitempty" json:"file_mime,omitempty"`
	FileExt            string    `bson:"file_ext,omitempty" json:"file_ext,omitempty"`
	IsCommit           bool      `bson:"is_commit,omitempty" json:"is_commit,omitempty"`
	IsPublic           bool      `bson:"is_public,omitempty" json:"is_public,omitempty"`
}

func (col *Media) ToJSON() common.GotaroMap {
	return common.GotaroMap{
		"id":                 col.ID,
		"created_at":         common.DateTimeNullableToString(&col.CreatedAt),
		"updated_at":         common.DateTimeNullableToString(&col.UpdatedAt),
		"deleted_at":         common.DateTimeNullableToString(&col.DeletedAt),
		"rule_slug":          col.RuleSlug,
		"driver_slug":        col.DriverSlug,
		"file_path":          col.FilePath,
		"file_original_name": col.FileOriginalName,
		"file_alias_name":    col.FileAliasName,
		"file_size":          col.FileSize,
		"file_mime":          col.FileMime,
		"file_ext":           col.FileExt,
		"is_commit":          col.IsCommit,
	}
}

func (col *Media) ToJSONSimple() common.GotaroMap {
	return common.GotaroMap{
		"id":                 col.ID,
		"rule_slug":          col.RuleSlug,
		"driver_slug":        col.DriverSlug,
		"file_path":          col.FilePath,
		"file_original_name": col.FileOriginalName,
		"file_alias_name":    col.FileAliasName,
		"file_size":          col.FileSize,
		"file_mime":          col.FileMime,
		"file_ext":           col.FileExt,
		"is_commit":          col.IsCommit,
		"is_public":          col.IsPublic,
	}
}

func (col *Media) ToMediaResult() common.GotaroMap {

	return common.GotaroMap{
		"id":               col.ID,
		"gotaro_file_path": col.GetGotaroFilePath(),
		"url":              col.FilePath,
		"is_public":        col.IsPublic,
	}
}

func (col *Media) GetGotaroFilePath() string {
	return fmt.Sprintf("gotaro://%s/%s", col.RuleSlug, col.FileAliasName)
}

func (col *Media) ToJSONString() (string, error) {
	data, err := json.Marshal(col)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (col *Media) FromJSONString(data string) error {
	return json.Unmarshal([]byte(data), col)
}

func (col Media) GetCollName() string {
	return "medias"
}
