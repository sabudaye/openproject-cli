package printer

import (
	"fmt"
	"strings"

	"github.com/sabudaye/openproject-cli/models"
)

func Relation(r *models.Relation) {
	activeRenderer.Relation(r)
}

func Relations(rs []*models.Relation) {
	activeRenderer.Relations(rs)
}

func printRelation(r *models.Relation, maxIdLength int) {
	diff := maxIdLength - idLength(r.Id)
	parts := []string{Red(fmt.Sprintf("%s#%d", indent(diff), r.Id))}

	if r.FromId > 0 {
		parts = append(parts, Red(fmt.Sprintf("#%d", r.FromId)))
	}
	if r.FromTitle != "" {
		parts = append(parts, Cyan(r.FromTitle))
	}

	parts = append(parts, "["+Yellow(r.Name)+"]")

	if r.ToId > 0 {
		parts = append(parts, Red(fmt.Sprintf("#%d", r.ToId)))
	}
	if r.ToTitle != "" {
		parts = append(parts, Cyan(r.ToTitle))
	}

	if r.HasLag {
		parts = append(parts, Green(fmt.Sprintf("(lag: %dd)", r.Lag)))
	}
	if r.Description != "" {
		parts = append(parts, Blue(fmt.Sprintf("%q", r.Description)))
	}

	activePrinter.Println(strings.Join(parts, " "))
}
