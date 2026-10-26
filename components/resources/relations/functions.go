package relations

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

// InvolvedFilter selects relations in which the given work package takes
// part on either side (from or to), per the "involved" relation filter
// of the APIv3 list_relations endpoint.
func InvolvedFilter(workPackageId string) requests.Filter {
	return requests.Filter{
		Operator: "=",
		Name:     "involved",
		Values:   []string{workPackageId},
	}
}

func ByWorkPackage(id string) ([]*models.Relation, error) {
	query := requests.NewUnpaginatedQuery(nil, []requests.Filter{InvolvedFilter(id)})

	response, err := requests.Get(paths.Relations(), &query)
	if err != nil {
		return nil, err
	}

	collection := parser.Parse[dtos.RelationCollectionDto](response)
	return collection.Convert(), nil
}
