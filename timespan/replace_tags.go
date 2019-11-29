package timespan

import (
	"context"

	"github.com/traggo/server/auth"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
)

// ReplaceTimeSpanTags replaces time span tags.
func (r *ResolverForTimeSpan) ReplaceTimeSpanTags(ctx context.Context, fromExternal gqlmodel.InputTimeSpanTag, toExternal gqlmodel.InputTimeSpanTag, opt gqlmodel.InputReplaceOptions) (*bool, error) {
	userID := auth.GetUser(ctx).ID

	from := tagToInternal(fromExternal)
	to := tagToInternal(toExternal)

	if err := tagsExist(r.DB, userID, []model.TimeSpanTag{from}); err != nil {
		return nil, err
	}
	if err := tagsExist(r.DB, userID, []model.TimeSpanTag{to}); err != nil {
		return nil, err
	}

	tx := r.DB.Begin()

	hasToKey := tx.Table("time_span_tags as innertst").
		Where("innertst.key = ?", to.Key).
		Where("innerts.id = innertst.time_span_id").
		SubQuery()

	if opt.Override == gqlmodel.OverrideModeOverride {
		timeSpanIdsWithExistingToKey := tx.Table("time_spans as innerts").
			Select("id").
			Where("innerts.user_id = ?", userID).
			Where("EXISTS ?", hasToKey).
			SubQuery()

		if err := tx.Where("time_span_id in ?", timeSpanIdsWithExistingToKey).
			Where(&model.TimeSpanTag{Key: to.Key}).
			Delete(new(model.TimeSpanTag)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	timeSpansIdsOfUser := tx.Table("time_spans as innerts").
		Select("id").
		Where("innerts.user_id = ?", userID).
		Where("NOT(EXISTS ?)", hasToKey).
		SubQuery()

	if update := tx.Model(&model.TimeSpanTag{}).
		Where("time_span_id in ?", timeSpansIdsOfUser).
		Where(from).
		Updates(to); update.Error != nil {
		tx.Rollback()
		return nil, update.Error
	}

	if opt.Override == gqlmodel.OverrideModeDiscard {
		if err := tx.Where(&from).Delete(new(model.TimeSpanTag)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	commit := tx.Commit()

	return nil, commit.Error
}
