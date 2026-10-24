package printer

import "github.com/sabudaye/openproject-cli/models"

func CustomActions(actions []*models.CustomAction) {
	activeRenderer.CustomActions(actions)
}
