package entry_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/dashboard"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestEntries(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	var bVal bool

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

	_, err = resolver.AddDashboardEntry(user1, 5, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    "",
		Tags:        []string{"hhol"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d",
		},
	}, nil)
	require.EqualError(t, err, "dashboard does not exist")
	_, err = resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalHourly,
		Tags:        []string{"hhol"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     p(55),
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now-2d",
		},
	}, nil)
	require.EqualError(t, err, "dashboard range does not exist")
	_, err = resolver.AddDashboardEntry(user2, 1, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    "doubly",
		Tags:        []string{"hhol"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "",
			To:   "",
		},
	}, nil)
	require.EqualError(t, err, "dashboard does not exist")
	_, err = resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"hhol"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2d",
			To:   "now-2",
		},
	}, nil)
	require.EqualError(t, err, "range to (now-2) invalid: expected unit at the end but got nothing")
	_, err = resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"hhol"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2",
			To:   "now-2d",
		},
	}, nil)
	require.EqualError(t, err, "range from (now-2) invalid: expected unit at the end but got nothing")
	_, err = resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2d",
			To:   "now-2d",
		},
	}, nil)
	require.EqualError(t, err, "at least one tag is required")

	entry, err := resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "test", false, gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"abc"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2d",
			To:   "now-5d",
		},
	}, &gqlmodel.InputResponsiveDashboardEntryPos{Desktop: &gqlmodel.InputDashboardEntryPos{
		H: 1,
		W: 2,
		X: 1,
		Y: 3,
	}})
	require.NoError(t, err)
	expectedEntry := &gqlmodel.DashboardEntry{
		ID:    1,
		Title: "test",
		Total: false,
		StatsSelection: &gqlmodel.StatsSelection{
			Interval:    gqlmodel.StatsIntervalDaily,
			Tags:        []string{"abc"},
			ExcludeTags: nil,
			IncludeTags: nil,
			RangeID:     nil,
			Range: &gqlmodel.RelativeOrStaticRange{
				From: "now-2d",
				To:   "now-5d",
			},
		},
		Pos: &gqlmodel.ResponsiveDashboardEntryPos{
			Desktop: &gqlmodel.DashboardEntryPos{
				W:    2,
				H:    3,
				X:    1,
				Y:    3,
				MinW: 2,
				MinH: 3,
			},
			Mobile: &gqlmodel.DashboardEntryPos{
				W:    1,
				H:    2,
				X:    0,
				Y:    0,
				MinW: 1,
				MinH: 2,
			},
		},
		EntryType: gqlmodel.EntryTypeBarChart,
	}
	require.Equal(t, expectedEntry, entry)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, dashboards[0].Items, []*gqlmodel.DashboardEntry{expectedEntry})
	xrange, err := resolver.AddDashboardRange(user1, 1, gqlmodel.InputNamedDateRange{
		Name:     "range",
		Editable: false,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-1d",
			To:   "now",
		},
	})
	require.NoError(t, err)
	_, err = resolver.AddDashboardEntry(user1, 1, gqlmodel.EntryTypeBarChart, "other", false, gqlmodel.InputStatsSelection{
		Interval: gqlmodel.StatsIntervalDaily,
		Tags:     []string{"abc"},
		RangeID:  p(xrange.ID),
	}, &gqlmodel.InputResponsiveDashboardEntryPos{Desktop: &gqlmodel.InputDashboardEntryPos{
		H: 1,
		W: 2,
		X: 1,
		Y: 3,
	}})
	require.NoError(t, err)
	_, err = resolver.RemoveDashboardRange(user1, xrange.ID)
	require.EqualError(t, err, "range is used in entries: other")

	chart := gqlmodel.EntryTypePieChart
	_, err = resolver.UpdateDashboardEntry(user2, 1, &chart, nil, nil, nil, nil)
	require.EqualError(t, err, "dashboard does not exist")
	_, err = resolver.UpdateDashboardEntry(user1, 3, &chart, nil, nil, nil, nil)
	require.EqualError(t, err, "entry does not exist")
	_, err = resolver.UpdateDashboardEntry(user1, 1, &chart, s("cool title"), &bVal, &gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"kek"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2",
			To:   "now-5d",
		},
	}, nil)
	require.EqualError(t, err, "range from (now-2) invalid: expected unit at the end but got nothing")
	_, err = resolver.UpdateDashboardEntry(user1, 1, &chart, s("cool title"), &bVal, &gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"kek"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2d",
			To:   "now-5",
		},
	}, nil)
	require.EqualError(t, err, "range to (now-5) invalid: expected unit at the end but got nothing")
	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.DashboardEntry{
		{
			ID:    1,
			Title: "test",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"abc"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     nil,
				Range: &gqlmodel.RelativeOrStaticRange{
					From: "now-2d",
					To:   "now-5d",
				},
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypeBarChart,
		},
		{
			ID:    2,
			Title: "other",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"abc"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     p(1),
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypeBarChart,
		},
	}, dashboards[0].Items)
	_, err = resolver.UpdateDashboardEntry(user1, 1, &chart, s("cool title"), &bVal, &gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"kek"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     nil,
		Range: &gqlmodel.InputRelativeOrStaticRange{
			From: "now-2d",
			To:   "now-12d",
		},
	}, &gqlmodel.InputResponsiveDashboardEntryPos{Desktop: &gqlmodel.InputDashboardEntryPos{
		H: 1,
		W: 2,
		X: 1,
		Y: 3,
	}})
	require.NoError(t, err)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.DashboardEntry{
		{
			ID:    1,
			Title: "cool title",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"kek"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     nil,
				Range: &gqlmodel.RelativeOrStaticRange{
					From: "now-2d",
					To:   "now-12d",
				},
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypePieChart,
		},
		{
			ID:    2,
			Title: "other",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"abc"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     p(1),
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypeBarChart,
		},
	}, dashboards[0].Items)

	_, err = resolver.UpdateDashboardEntry(user1, 1, &chart, s("cool title"), &bVal, &gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"kek"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     p(5),
	}, &gqlmodel.InputResponsiveDashboardEntryPos{Desktop: &gqlmodel.InputDashboardEntryPos{
		H: 1,
		W: 2,
		X: 1,
		Y: 3,
	}})
	require.EqualError(t, err, "dashboard range does not exist")

	_, err = resolver.UpdateDashboardEntry(user1, 1, &chart, s("cool title"), &bVal, &gqlmodel.InputStatsSelection{
		Interval:    gqlmodel.StatsIntervalDaily,
		Tags:        []string{"kek"},
		ExcludeTags: nil,
		IncludeTags: nil,
		RangeID:     p(xrange.ID),
	}, &gqlmodel.InputResponsiveDashboardEntryPos{Desktop: &gqlmodel.InputDashboardEntryPos{
		H: 1,
		W: 2,
		X: 1,
		Y: 3,
	}})
	require.NoError(t, err)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.DashboardEntry{
		{
			ID:    1,
			Title: "cool title",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"kek"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     p(xrange.ID),
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypePieChart,
		},
		{
			ID:    2,
			Title: "other",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"abc"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     p(1),
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypeBarChart,
		},
	}, dashboards[0].Items)

	_, err = resolver.RemoveDashboardEntry(user2, 1)
	require.EqualError(t, err, "dashboard does not exist")

	_, err = resolver.RemoveDashboardEntry(user1, 55)
	require.EqualError(t, err, "entry does not exist")

	_, err = resolver.RemoveDashboardEntry(user1, 1)
	require.NoError(t, err)

	dashboards, err = resolver.Dashboards(user1)
	require.NoError(t, err)
	require.Equal(t, []*gqlmodel.DashboardEntry{
		{
			ID:    2,
			Title: "other",
			Total: false,
			StatsSelection: &gqlmodel.StatsSelection{
				Interval:    gqlmodel.StatsIntervalDaily,
				Tags:        []string{"abc"},
				ExcludeTags: nil,
				IncludeTags: nil,
				RangeID:     p(1),
			},
			Pos: &gqlmodel.ResponsiveDashboardEntryPos{
				Desktop: &gqlmodel.DashboardEntryPos{
					W:    2,
					H:    3,
					X:    1,
					Y:    3,
					MinW: 2,
					MinH: 3,
				},
				Mobile: &gqlmodel.DashboardEntryPos{
					W:    1,
					H:    2,
					X:    0,
					Y:    0,
					MinW: 1,
					MinH: 2,
				},
			},
			EntryType: gqlmodel.EntryTypeBarChart,
		},
	}, dashboards[0].Items)
}

func p(i int) *int {
	return &i
}
func s(s string) *string {
	return &s
}
