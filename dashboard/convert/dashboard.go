package convert

import (
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// ToExternalDashboards converts dashboards.
func ToExternalDashboards(dashboards []model.Dashboard) ([]*gqlmodel.Dashboard, error) {
	result := []*gqlmodel.Dashboard{}
	for _, dashboard := range dashboards {
		if converted, err := ToExternalDashboard(dashboard); err == nil {
			result = append(result, converted)
		} else {
			return nil, err
		}
	}
	return result, nil
}

// ToExternalDashboard converts a dashboard.
func ToExternalDashboard(dashboard model.Dashboard) (*gqlmodel.Dashboard, error) {
	entries, err := toExternalEntries(dashboard.Entries)
	if err != nil {
		return nil, err
	}
	ranges, err := toExternalDashboardsRanges(dashboard.Ranges)
	if err != nil {
		return nil, err
	}
	return &gqlmodel.Dashboard{
		ID:     dashboard.ID,
		Name:   dashboard.Name,
		Items:  entries,
		Ranges: ranges,
	}, nil
}
