package types

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

func All() ([]*models.Type, error) {
	query := requests.NewPaginatedQuery(-1, nil)
	response, err := requests.Get(paths.Types(), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.TypeCollectionDto](response)
	return element.Convert(), nil
}
