package notifications

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

func All(reason string) ([]*models.Notification, error) {
	response, err := requests.Get(paths.Notifications(), generateQuery(reason))
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.NotificationCollectionDto](response)
	return element.Convert(), nil
}

func generateQuery(reason string) *requests.Query {
	filters := []requests.Filter{
		ReadFilter(false),
	}

	if reason != "" {
		filters = append(filters, ReasonFilter(reason))
	}

	query := requests.NewFilterQuery(filters)
	return &query
}
