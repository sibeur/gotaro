package dto

type AuthDTO struct {
	APIKey    string `json:"api_key" validate:"required"`
	SecretKey string `json:"secret_key" validate:"required"`
}
