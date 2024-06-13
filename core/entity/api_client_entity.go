package entity

import (
	"time"

	"github.com/sibeur/gotaro/core/common"
)

type APIClient struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Name      string    `bson:"name,omitempty" json:"name,omitempty"`
	Key       string    `bson:"key,omitempty" json:"key,omitempty"`
	Secret    string    `bson:"secret,omitempty" json:"secret,omitempty"`
	Scopes    []string  `bson:"scopes,omitempty" json:"scopes,omitempty"`
}

func (a *APIClient) ToJSON() common.GotaroMap {
	return common.GotaroMap{
		"id":         a.ID,
		"name":       a.Name,
		"key":        a.Key,
		"secret":     a.Secret,
		"scopes":     a.Scopes,
		"created_at": a.CreatedAt,
		"updated_at": a.UpdatedAt,
	}
}

func (a *APIClient) ToJSONSimple() common.GotaroMap {
	return common.GotaroMap{
		"id":         a.ID,
		"name":       a.Name,
		"created_at": a.CreatedAt,
		"updated_at": a.UpdatedAt,
	}
}

func (a APIClient) GetCollName() string {
	return "api_clients"
}
