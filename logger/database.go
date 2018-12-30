package logger

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	sqlRegexp = regexp.MustCompile(`\?`)
)

// DatabaseLogger logs sql queries
type DatabaseLogger struct {
}

// Print pretty prints the gorm.DB log values
// Mostly copied from https://github.com/jinzhu/gorm/blob/master/logger.go
func (l *DatabaseLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
		)

		if level == "sql" {
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if _, ok := value.([]byte); ok {
						formattedValues = append(formattedValues, "'<binary>'")
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			formattedValuesLength := len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}

			log.Debug().Str("took", values[2].(time.Duration).String()).Int64("rows", values[5].(int64)).Msg("SQL: " + sql)
		} else if level == "log" {
			if len(values) > 2 {
				if err, ok := values[2].(error); ok {
					log.Error().Err(err).Msg("Database error")
					return
				}
			}
			log.Debug().Msg(fmt.Sprint(values[2:]...))
		} else {
			log.Error().Msgf("DatabaseLogger: cannot handle level %#v", level)
		}
	}
	// ignore empty log
}
