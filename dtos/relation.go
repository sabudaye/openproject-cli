package dtos

import (
	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/models"
)

// RelationEmbeddedWorkPackageDto mirrors the abbreviated work package
// embedded in a relation when the instance returns _embedded.from/to.
type RelationEmbeddedWorkPackageDto struct {
	Id int64 `json:"id,omitempty"`
}

type RelationEmbeddedDto struct {
	From *RelationEmbeddedWorkPackageDto `json:"from,omitempty"`
	To   *RelationEmbeddedWorkPackageDto `json:"to,omitempty"`
}

type RelationLinksDto struct {
	Self *LinkDto `json:"self,omitempty"`
	From *LinkDto `json:"from,omitempty"`
	To   *LinkDto `json:"to,omitempty"`
}

type RelationDto struct {
	Id          int64                `json:"id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Type        string               `json:"type,omitempty"`
	ReverseType string               `json:"reverseType,omitempty"`
	Lag         *int                 `json:"lag,omitempty"`
	Description string               `json:"description,omitempty"`
	Embedded    *RelationEmbeddedDto `json:"_embedded,omitempty"`
	Links       *RelationLinksDto    `json:"_links,omitempty"`
}

type relationElements struct {
	Elements []*RelationDto `json:"elements"`
}

type RelationCollectionDto struct {
	Embedded relationElements `json:"_embedded"`
	Type     string           `json:"_type"`
	Total    int64            `json:"total"`
	Count    int64            `json:"count"`
}

/////////////// MODEL CONVERSION ///////////////

// wpIdFromRef resolves a work package id from a relation reference,
// preferring an explicit embedded id and falling back to parsing the
// _links href (e.g. /api/v3/work_packages/108 -> 108).
func wpIdFromRef(emb *RelationEmbeddedWorkPackageDto, link *LinkDto) uint64 {
	if emb != nil && emb.Id > 0 {
		return uint64(emb.Id)
	}
	if link != nil && link.Href != "" {
		return parser.IdFromLink(link.Href)
	}
	return 0
}

func (dto *RelationDto) Convert() *models.Relation {
	var lag int
	var hasLag bool
	if dto.Lag != nil {
		lag = *dto.Lag
		hasLag = true
	}

	var embFrom, embTo *RelationEmbeddedWorkPackageDto
	var linkFrom, linkTo *LinkDto
	if dto.Embedded != nil {
		embFrom = dto.Embedded.From
		embTo = dto.Embedded.To
	}
	if dto.Links != nil {
		linkFrom = dto.Links.From
		linkTo = dto.Links.To
	}

	return &models.Relation{
		Id:          uint64(dto.Id),
		Name:        dto.Name,
		Type:        dto.Type,
		ReverseType: dto.ReverseType,
		Lag:         lag,
		HasLag:      hasLag,
		Description: dto.Description,
		FromId:      wpIdFromRef(embFrom, linkFrom),
		FromTitle:   linkTitle(linkFrom),
		ToId:        wpIdFromRef(embTo, linkTo),
		ToTitle:     linkTitle(linkTo),
	}
}

func linkTitle(link *LinkDto) string {
	if link == nil {
		return ""
	}
	return link.Title
}

func (dto *RelationCollectionDto) Convert() []*models.Relation {
	var relations = make([]*models.Relation, len(dto.Embedded.Elements))

	for idx, r := range dto.Embedded.Elements {
		relations[idx] = r.Convert()
	}

	return relations
}
