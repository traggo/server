package dbrange_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/dashboard"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestRanges(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := dashboard.NewResolverForDashboard(db.DB)

	// initial no dashboards
	user1 := fake.User(1)
	user2 := fake.User(2)
	db.User(1)
	db.User(2)
	dashboards, err := resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Empty(t, dashboards)

	dashboard, err := resolver.CreateDashboard(user1, "cool dashboard")
	require.NoError(t, err)
	expectAdded := &gqlmodel.Dashboard{
		ID:     1,
		Name:   "cool dashboard",
		Ranges: []*gqlmodel.NamedDateRange{},
		Items:  []*gqlmodel.DashboardEntry{},
	}
	require.Equal(t, expectAdded, dashboard)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Dashboard{expectAdded}, dashboards)

	dashboards, err = resolver.Dashboards(user2)
	require.NoError(t, err)
	require.Empty(t, dashboards)

	xrange, err := resolver.AddDashboardRange(user1, dashboard.ID, gqlmodel.InputNamedDateRange{
		Name:     "new range",
		Editable: false,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d",
		},
	})
	require.NoError(t, err)
	require.Equal(t, &gqlmodel.NamedDateRange{
		ID:       1,
		Name:     "new range",
		Editable: false,
		Range: &gqlmodel.RelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d",
		},
	}, xrange)
	_, err = resolver.AddDashboardRange(user1, dashboard.ID, gqlmodel.InputNamedDateRange{
		Name:     "new range",
		Editable: false,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1",
			To:   "now-2d",
		},
	})
	require.EqualError(t, err, "range from (now-1) invalid: expected unit at the end but got nothing")
	_, err = resolver.AddDashboardRange(user1, dashboard.ID, gqlmodel.InputNamedDateRange{
		Name:     "new range",
		Editable: false,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2",
		},
	})
	require.EqualError(t, err, "range to (now-2) invalid: expected unit at the end but got nothing")
	_, err = resolver.AddDashboardRange(user1, 100, gqlmodel.InputNamedDateRange{
		Name:     "new range",
		Editable: false,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2",
		},
	})
	require.EqualError(t, err, "dashboard does not exist")
	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Dashboard{{
		ID:   1,
		Name: "cool dashboard",
		Ranges: []*gqlmodel.NamedDateRange{{
			ID:       1,
			Name:     "new range",
			Editable: false,
			Range: &gqlmodel.RelativeOrStaticRange{
				From: "now-1d",
				To:   "now-2d",
			},
		}},
		Items: []*gqlmodel.DashboardEntry{},
	}}, dashboards)

	_, err = resolver.UpdateDashboardRange(user2, 1, gqlmodel.InputNamedDateRange{
		Name:     "my range",
		Editable: true,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d/w",
		},
	})
	require.EqualError(t, err, "dashboard does not exist")
	xrange, err = resolver.UpdateDashboardRange(user1, 1, gqlmodel.InputNamedDateRange{
		Name:     "my range",
		Editable: true,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d/w",
		},
	})
	require.Equal(t, &gqlmodel.NamedDateRange{
		ID:       1,
		Name:     "my range",
		Editable: true,
		Range: &gqlmodel.RelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d/w",
		},
	}, xrange)
	_, err = resolver.UpdateDashboardRange(user1, 1, gqlmodel.InputNamedDateRange{
		Name:     "my range",
		Editable: true,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2",
		},
	})
	require.EqualError(t, err, "range to (now-2) invalid: expected unit at the end but got nothing")
	_, err = resolver.UpdateDashboardRange(user1, 44, gqlmodel.InputNamedDateRange{
		Name:     "my range",
		Editable: true,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2",
		},
	})
	require.EqualError(t, err, "dashboard range does not exist")

	entry, err := resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "other", true, gqlmodel.InputStatsSelection{
		Interval: gqlmodel.StatsIntervalDaily,
		Tags:     []string{"abc"},
		RangeID:  &xrange.ID,
	}, &gqlmodel.InputResponsiveDashboardEntryPos{Desktop: &gqlmodel.InputDashboardEntryPos{
		H: 1,
		W: 2,
		X: 1,
		Y: 3,
	}})
	require.NoError(t, err)
	_, err = resolver.RemoveDashboardRange(user1, xrange.ID)
	require.EqualError(t, err, "range is used in entries: other")
	_, err = resolver.RemoveDashboardEntry(user1, entry.ID)
	require.NoError(t, err)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.Dashboard{{
		ID:   1,
		Name: "cool dashboard",
		Ranges: []*gqlmodel.NamedDateRange{{
			ID:       1,
			Name:     "my range",
			Editable: true,
			Range: &gqlmodel.RelativeOrStaticRange{
				From: "now-1d",
				To:   "now-2d/w",
			},
		}},
		Items: []*gqlmodel.DashboardEntry{},
	}}, dashboards)

	_, err = resolver.RemoveDashboardRange(user2, 1)
	require.EqualError(t, err, "dashboard does not exist")
	_, err = resolver.RemoveDashboardRange(user2, 55)
	require.EqualError(t, err, "dashboard range does not exist")

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Len(t, dashboards[0].Ranges, 1)

	_, err = resolver.RemoveDashboardRange(user1, 1)
	require.NoError(t, err)
	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Len(t, dashboards[0].Ranges, 0)

}
