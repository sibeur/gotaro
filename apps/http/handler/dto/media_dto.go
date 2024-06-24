package dto

type GetMediaBatchDTO struct {
	Files []string `json:"files" validate:"required"`
}
