package main

import (
	"flag"
	"fmt"
	"math/rand"
	gotime "time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/database"
	"github.com/traggo/server/logger"
	"github.com/traggo/server/model"
	"github.com/traggo/server/time"
	"github.com/traggo/server/user/password"
)

func main() {
	logger.Init(zerolog.InfoLevel)
	rand.Seed(0)
	var from string
	flag.StringVar(&from, "from", "now/y", "from")
	var to string
	flag.StringVar(&to, "to", "now", "to")

	db, err := database.New("sqlite3", "data.db")
	noErr(err)
	defer db.Close()

	log.Info().Msg("Creating User ...")
	user := &model.User{
		Name:  "admin",
		Pass:  password.CreatePassword("admin", 10),
		Admin: true,
	}
	noErr(db.Create(user).Error)
	log.Info().Msg("... Done")
	uID := user.ID

	log.Info().Msg("Creating Tags ...")

	db.Create(&model.TagDefinition{
		Type:   model.TypeSingleValue,
		Key:    "type",
		UserID: uID,
		Color:  "#fff",
	})
	db.Create(&model.TagDefinition{
		Type:   model.TypeSingleValue,
		Key:    "issue",
		UserID: uID,
		Color:  "#fff",
	})
	db.Create(&model.TagDefinition{
		Type:   model.TypeSingleValue,
		Key:    "meeting",
		UserID: uID,
		Color:  "#fff",
	})
	db.Create(&model.TagDefinition{
		Type:   model.TypeSingleValue,
		Key:    "with",
		UserID: uID,
		Color:  "#fff",
	})
	db.Create(&model.TagDefinition{
		Type:   model.TypeSingleValue,
		Key:    "misc",
		UserID: uID,
		Color:  "#fff",
	})
	db.Create(&model.TagDefinition{
		Type:   model.TypeSingleValue,
		Key:    "proj",
		UserID: uID,
		Color:  "#fff",
	})
	log.Info().Msg("... Done")

	now := gotime.Now().UTC()
	fromTime, err := time.ParseTime(now, from, true, gotime.Monday)
	noErr(err)
	toTime, err := time.ParseTime(now, to, false, gotime.Sunday)
	noErr(err)

	buckets := []bucket{
		{
			TimeSpans: generateSupport(),
			weight:    100,
		},
		{
			TimeSpans: generateMeeting(),
			weight:    30,
		},
		{
			TimeSpans: generateIssueType("gotify", "review"),
			weight:    100,
			Active:    6,
		},
		{
			TimeSpans: generateIssueType("gotify", "work"),
			weight:    100,
			Active:    6,
		},
		{
			TimeSpans: generateIssueType("traggo", "review"),
			weight:    100,
			Active:    6,
		},
		{
			TimeSpans: generateIssueType("traggo", "work"),
			weight:    100,
			Active:    6,
		},
	}

	sum := 0
	for _, bucket := range buckets {
		sum += bucket.weight
	}

	currentNow := fromTime

	log.Info().Msg("Creating TimeSpans ...")
	for i := 1; currentNow.Before(toTime); i++ {
		selected := randInt(0, sum)
		current := 0
		var bucket bucket
		for _, bucket = range buckets {
			current += bucket.weight
			if current > selected {
				break
			}
		}

		duration := gotime.Minute * gotime.Duration(randInt(1, 150))

		ts := bucket.TimeSpans[randInt(0, len(bucket.TimeSpans))]
		end := currentNow.Add(duration)
		db.Create(&model.TimeSpan{
			Tags:          ts.Tags,
			UserID:        uID,
			StartUserTime: currentNow,
			StartUTC:      currentNow,
			EndUserTime:   &end,
			EndUTC:        &end,
			OffsetUTC:     0,
		})
		currentNow = end.Add(gotime.Minute * gotime.Duration(randInt(1, 20)))
		if i%500 == 0 {
			log.Info().Msgf("Inserted %d ...", i)
		}
	}
	log.Info().Msg("... Done")
}

func generateIssueType(proj, t string) []timeSpan {
	numbers := []int{}
	for i := 1; i <= 1000; i++ {
		numbers = append(numbers, i)
	}

	result := []timeSpan{}
	for _, number := range numbers {
		result = append(result, timeSpan{
			Tags: []model.TimeSpanTag{
				{Key: "proj", StringValue: &proj},
				{Key: "type", StringValue: p(t)},
				{Key: "issue", StringValue: p(fmt.Sprintf("%X-%d", proj, number))},
			},
			Runtime: gotime.Hour * 24 * 12,
		})
	}
	return result
}

func generateMeeting() []timeSpan {
	names := []string{"weekly", "retro", "company", "team", "feature"}

	result := []timeSpan{}
	for _, name := range names {
		result = append(result, timeSpan{
			Tags: []model.TimeSpanTag{
				{Key: "type", StringValue: p("meeting")},
				{Key: "meeting", StringValue: &name},
			},
			Runtime: gotime.Duration(-1),
		})
	}
	return result
}

func generateSupport() []timeSpan {
	names := []string{"Tom", "Jerry", "Nico", "Anton", "Noob", "Alex", "Bob", "jmattheis", "nicories"}

	result := []timeSpan{}
	for _, name := range names {
		result = append(result, timeSpan{
			Tags: []model.TimeSpanTag{
				{Key: "type", StringValue: p("support")},
				{Key: "with", StringValue: &name},
			},
			Runtime: gotime.Duration(-1),
		})
	}
	return result
}

type bucket struct {
	TimeSpans       []timeSpan
	Active          int
	CurrentlyActive []timeSpan
	weight          int
}

type timeSpan struct {
	ActiveSince gotime.Time
	Runtime     gotime.Duration
	Tags        []model.TimeSpanTag
}

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func p(s string) *string {
	return &s
}
