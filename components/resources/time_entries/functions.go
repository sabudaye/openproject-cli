package time_entries

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

func All(query requests.Query) ([]*models.TimeEntry, error) {
	response, err := requests.Get(paths.TimeEntries(), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.TimeEntryCollectionDto](response)
	return element.Convert(), nil
}
