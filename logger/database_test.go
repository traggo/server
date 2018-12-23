package logger_test

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/traggo/server/logger"
	"github.com/traggo/server/test"
)

type Test struct {
	Password    []byte
	NullValue   *string
	Placeholder string
	Date        time.Time
	Number      uint
}

func TestDatabaseLogger_Success(t *testing.T) {
	someTime, _ := time.Parse("2006/01/02", "2017/01/02")
	db := test.InMemoryDB(t)
	defer db.Close()
	db.AutoMigrate(&Test{})

	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	db.Create(&Test{
		Password:    []byte{1},
		Date:        someTime,
		NullValue:   nil,
		Number:      1,
		Placeholder: "hello there",
	})

	expected := `SQL: INSERT INTO "tests" ("password","null_value","placeholder","date","number") VALUES ('<binary>',NULL,'hello there','2017-01-02 00:00:00','1')`

	fakeLog.AssertCount(1)
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.DebugLevel, Message: expected})
}

func TestDatabaseLogger_Error(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.AutoMigrate(&Test{})

	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	// column id doesn't exists
	db.Model(new(Test)).Where("id = ?", 5).Update("number", 6)

	expectedOne := `SQL: UPDATE "tests" SET "number" = '6'  WHERE (id = '5')`
	expectedTwo := `Database error`

	fakeLog.AssertCount(2)
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.DebugLevel, Message: expectedOne})
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.ErrorLevel, Message: expectedTwo})
}

func TestDatabaseLogger_EmptyLog(t *testing.T) {
	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	databaseLogger := logger.DatabaseLogger{}
	databaseLogger.Print([]interface{}{})

	fakeLog.AssertCount(0)
}

func TestDatabaseLogger_UnkownType(t *testing.T) {
	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	databaseLogger := logger.DatabaseLogger{}
	databaseLogger.Print("abc", "somevalue")

	expected := `DatabaseLogger: cannot handle level "abc"`

	fakeLog.AssertCount(1)
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.ErrorLevel, Message: expected})
}

func TestDatabaseLogger_LogWithoutError(t *testing.T) {
	fakeLog := test.NewLogger(t)
	defer fakeLog.Dispose()

	databaseLogger := logger.DatabaseLogger{}
	databaseLogger.Print("log", "somefile", "my log message")

	expected := `my log message`

	fakeLog.AssertCount(1)
	fakeLog.AssertEntryExists(test.Entry{Level: zerolog.DebugLevel, Message: expected})
}
