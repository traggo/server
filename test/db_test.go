package test_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestInMemoryDB(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	assert.NotNil(t, db)
}

func TestDatabase(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(1)

	user.AssertHasDevice("test device", false)
	user.NewDevice(1, "lol", "test device")
	user.AssertHasDevice("test device", true)

	ts := user.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")

	ts.AssertHasTagIgnoreValue("abc", false)
	ts.AssertHasTag("abc", "def", false)
	ts.Tag("abc", "def")
	ts.AssertHasTag("abc", "def", true)
	ts.AssertHasTagIgnoreValue("abc", true)

	db.NewUser(2, "abc", true)
	db.NewUserPass(3, "xxx", []byte{5, 5}, true)

	expectedUsers := []model.User{
		{ID: 1, Name: "test1", Pass: []uint8{1, 2, 3}, Admin: false},
		{ID: 2, Name: "abc", Pass: []uint8{1, 2, 3}, Admin: true},
		{ID: 3, Name: "xxx", Pass: []uint8{5, 5}, Admin: true}}

	var users []model.User
	db.Find(&users)
	assert.Equal(t, expectedUsers, users)

	expectedDevices := []model.Device{{
		ID:        1,
		Token:     "lol",
		Name:      "test device",
		UserID:    1,
		Type:      model.TypeNoExpiry,
		ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
		CreatedAt: test.Time("2009-06-30T18:30:00Z")}}

	var devices []model.Device
	db.Find(&devices)
	assert.Equal(t, expectedDevices, devices)

	expectedTimeSpans := []model.TimeSpan{{
		ID:            1,
		UserID:        1,
		StartUserTime: test.Time("2009-06-30T18:30:00Z"),
		StartUTC:      test.Time("2009-06-30T18:30:00Z"),
		EndUserTime:   test.TimeP("2009-06-30T18:40:00Z"),
		EndUTC:        test.TimeP("2009-06-30T18:40:00Z"),
		Tags: []model.TimeSpanTag{{
			Key:         "abc",
			StringValue: "def",
			TimeSpanID:  1,
		}},
	}}

	var timeSpans []model.TimeSpan
	db.Preload("Tags").Find(&timeSpans)
	assert.Equal(t, expectedTimeSpans, timeSpans)

	ts.AssertExists(true)
	db.Delete(new(model.TimeSpan), "id = ?", ts.TimeSpan.ID)
	ts.AssertExists(false)

	user.AssertHasTagDefinition("oops", false)
	user.NewTagDefinition("oops")
	user.AssertHasTagDefinition("oops", true)

	dash := user.Dashboard("hello")

	dash.AssertHasRange("hello", false)
	dash.Range("hello")
	dash.AssertHasRange("hello", true)

	dash.AssertHasEntry("abc", false)
	dash.Entry("abc")
	dash.AssertHasEntry("abc", true)

	dash.AssertExists(true)
	db.Delete(new(model.Dashboard), "id = ?", dash.Dashboard.ID)
	dash.AssertExists(false)

	user.AssertExists(true)
	db.Delete(new(model.User), "id = ?", user.User.ID)
	user.AssertExists(false)

}
