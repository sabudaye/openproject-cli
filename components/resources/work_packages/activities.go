package work_packages

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

func Activities(id string) (activites []*models.Activity, err error) {
	response, err := requests.Get(paths.WorkPackageActivities(id), nil)
	if err != nil {
		return nil, err
	}

	activitiesDto := parser.Parse[dtos.ActivityCollectionDto](response)
	return activitiesDto.Convert()
}
