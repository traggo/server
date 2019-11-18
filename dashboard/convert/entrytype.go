package convert

import (
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// InternalEntryType converts entry type.
func InternalEntryType(entryType gqlmodel.EntryType) model.DashboardType {
	switch entryType {
	case gqlmodel.EntryTypeBarChart:
		return model.TypeBarChart
	case gqlmodel.EntryTypePieChart:
		return model.TypePieChart
	case gqlmodel.EntryTypeStackedBarChart:
		return model.TypeStackedBarChart
	case gqlmodel.EntryTypeLineChart:
		return model.TypeLineChart
	case gqlmodel.EntryTypeHorizontalTable:
		return model.HorizontalTable
	case gqlmodel.EntryTypeVerticalTable:
		return model.VerticalTable
	default:
		panic("unknown entry type " + entryType)
	}
}

// ExternalEntryType converts entry type.
func ExternalEntryType(entryType model.DashboardType) gqlmodel.EntryType {
	switch entryType {
	case model.TypeBarChart:
		return gqlmodel.EntryTypeBarChart
	case model.TypePieChart:
		return gqlmodel.EntryTypePieChart
	case model.TypeStackedBarChart:
		return gqlmodel.EntryTypeStackedBarChart
	case model.TypeLineChart:
		return gqlmodel.EntryTypeLineChart
	case model.HorizontalTable:
		return gqlmodel.EntryTypeHorizontalTable
	case model.VerticalTable:
		return gqlmodel.EntryTypeVerticalTable
	default:
		panic("unknown entry type " + entryType)
	}
}
