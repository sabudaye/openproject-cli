package dtos

import "github.com/sabudaye/openproject-cli/models"

type LongTextDto struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	Html   string `json:"html"`
}

// ///////////// MODEL CONVERSION ///////////////
func (dto *LongTextDto) Convert() *models.LongText {
	return &models.LongText{
		Format: dto.Format,
		Raw:    dto.Raw,
		Html:   dto.Html,
	}
}
