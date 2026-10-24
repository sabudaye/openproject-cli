package projects

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

func AvailableVersions(projectId string) ([]*models.Version, error) {
	response, err := requests.Get(paths.ProjectVersions(projectId), nil)
	if err != nil {
		return nil, err
	}

	versionCollection := parser.Parse[dtos.VersionCollectionDto](response)

	return versionCollection.Convert(), nil
}
