package convert

import (
	"fmt"

	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/time"
)

func toExternalDashboardsRanges(ranges []model.DashboardRange) ([]*gqlmodel.NamedDateRange, error) {
	result := []*gqlmodel.NamedDateRange{}
	for _, xrange := range ranges {
		result = append(result, ToExternalDashboardRange(xrange))
	}
	return result, nil
}

// ToExternalDashboardRange converts a range.
func ToExternalDashboardRange(xrange model.DashboardRange) *gqlmodel.NamedDateRange {
	return &gqlmodel.NamedDateRange{
		ID:       xrange.ID,
		Name:     xrange.Name,
		Editable: xrange.Editable,
		Range: &gqlmodel.RelativeOrStaticRange{
			From: xrange.From,
			To:   xrange.To,
		},
	}
}

// ToInternalDashboardRange converts a range.
func ToInternalDashboardRange(xrange gqlmodel.InputNamedDateRange) (model.DashboardRange, error) {
	if err := time.Validate(xrange.Range.From); err != nil {
		return model.DashboardRange{}, fmt.Errorf("range from (%s) invalid: %s", xrange.Range.From, err)
	}
	if err := time.Validate(xrange.Range.To); err != nil {
		return model.DashboardRange{}, fmt.Errorf("range to (%s) invalid: %s", xrange.Range.To, err)
	}
	return model.DashboardRange{
		Editable: xrange.Editable,
		From:     xrange.Range.From,
		To:       xrange.Range.To,
		Name:     xrange.Name,
	}, nil
}
