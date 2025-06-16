package convert

import (
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

func excludedTagsToExternal(tags []model.DashboardExcludedTag) []*gqlmodel.TimeSpanTag {
	result := []*gqlmodel.TimeSpanTag{}
	for _, tag := range tags {
		result = append(result, &gqlmodel.TimeSpanTag{
			Key:   tag.Key,
			Value: tag.StringValue,
		})
	}
	return result
}

func ExcludedTagsToInternal(gqls []*gqlmodel.InputTimeSpanTag) []model.DashboardExcludedTag {
	result := make([]model.DashboardExcludedTag, 0)
	for _, tag := range gqls {
		result = append(result, excludedTagToInternal(*tag))
	}
	return result
}

func excludedTagToInternal(gqls gqlmodel.InputTimeSpanTag) model.DashboardExcludedTag {
	return model.DashboardExcludedTag{
		Key:         gqls.Key,
		StringValue: gqls.Value,
	}
}

func includedTagsToExternal(tags []model.DashboardIncludedTag) []*gqlmodel.TimeSpanTag {
	result := []*gqlmodel.TimeSpanTag{}
	for _, tag := range tags {
		result = append(result, &gqlmodel.TimeSpanTag{
			Key:   tag.Key,
			Value: tag.StringValue,
		})
	}
	return result
}

func IncludedTagsToInternal(gqls []*gqlmodel.InputTimeSpanTag) []model.DashboardIncludedTag {
	result := make([]model.DashboardIncludedTag, 0)
	for _, tag := range gqls {
		result = append(result, includedTagToInternal(*tag))
	}
	return result
}

func includedTagToInternal(gqls gqlmodel.InputTimeSpanTag) model.DashboardIncludedTag {
	return model.DashboardIncludedTag{
		Key:         gqls.Key,
		StringValue: gqls.Value,
	}
}
