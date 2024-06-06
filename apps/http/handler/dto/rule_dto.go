package dto

type NewRuleDTO struct {
	Slug       string   `json:"slug" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	MaxSize    uint64   `json:"max_size"`
	Mimes      []string `json:"mimes"`
	DriverSlug string   `json:"driver_slug" validate:"required"`
}

type EditRuleDTO struct {
	Name       string   `json:"name" validate:"required"`
	MaxSize    uint64   `json:"max_size"`
	Mimes      []string `json:"mimes"`
	DriverSlug string   `json:"driver_slug" validate:"required"`
}
