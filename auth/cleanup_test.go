package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestCleanUp_stops(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	methodDone := make(chan struct{})
	stopCleanUp := make(chan bool)

	go func() {
		CleanUp(db.DB, time.Microsecond*10, stopCleanUp)
		methodDone <- struct{}{}
	}()

	stopCleanUp <- true

	select {
	case <-methodDone:
		// expected
	case <-time.After(time.Second):
		t.Fail()
	}
}

func TestCleanUp_removeExpiredDevices(t *testing.T) {
	now := test.Time("2018-06-30T18:30:00Z")
	timeDispose := fakeTime(now)
	defer timeDispose()

	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(2)
	db.Create(&model.Device{
		ID:        1,
		Token:     "abc",
		UserID:    2,
		Name:      "android 1",
		ExpiresAt: test.Time("2018-07-30T18:30:01Z"),
		ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
		CreatedAt: test.Time("2009-06-30T18:30:00Z"),
	})
	db.Create(&model.Device{
		ID:        2,
		Token:     "abc2",
		UserID:    2,
		Name:      "android 2",
		ExpiresAt: test.Time("2018-06-29T18:30:00Z"),
		ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
		CreatedAt: test.Time("2009-06-30T18:30:00Z"),
	})
	db.Create(&model.Device{
		ID:        3,
		Token:     "abc3",
		UserID:    2,
		Name:      "android 3",
		ExpiresAt: test.Time("2009-06-29T18:30:00Z"),
		ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
		CreatedAt: test.Time("2009-06-30T18:30:00Z"),
	})
	db.Create(&model.Device{
		ID:        4,
		Token:     "abc4",
		UserID:    2,
		Name:      "android 4",
		ExpiresAt: test.Time("2025-06-29T18:30:00Z"),
		ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
		CreatedAt: test.Time("2009-06-30T18:30:00Z"),
	})

	stopCleanUp := make(chan bool)

	go CleanUp(db.DB, time.Microsecond*10, stopCleanUp)

	<-time.After(time.Millisecond * 20)

	expected := []model.Device{
		{
			ID:        1,
			Token:     "abc",
			UserID:    2,
			Name:      "android 1",
			ExpiresAt: test.Time("2018-07-30T18:30:01Z"),
			ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
			CreatedAt: test.Time("2009-06-30T18:30:00Z"),
		},
		{
			ID:        4,
			Token:     "abc4",
			UserID:    2,
			Name:      "android 4",
			ExpiresAt: test.Time("2025-06-29T18:30:00Z"),
			ActiveAt:  test.Time("2009-06-30T18:30:00Z"),
			CreatedAt: test.Time("2009-06-30T18:30:00Z"),
		},
	}

	var devices []model.Device
	db.Find(&devices)
	assert.Equal(t, expected, devices)
	stopCleanUp <- true
}

func fakeTime(t time.Time) func() {
	old := timeNow
	timeNow = func() time.Time {
		return t
	}
	return func() {
		timeNow = old
	}
}
