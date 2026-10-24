package printer

import "github.com/sabudaye/openproject-cli/models"

func Budget(budget *models.Budget) {
	activeRenderer.Budget(budget)
}

func Budgets(budgets []*models.Budget) {
	activeRenderer.Budgets(budgets)
}
