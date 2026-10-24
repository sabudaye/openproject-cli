package printer

import "github.com/sabudaye/openproject-cli/models"

func Status(s *models.Status) {
	activeRenderer.Status(s)
}

func StatusList(statuses []*models.Status) {
	activeRenderer.StatusList(statuses)
}
