package db

import (
	"runtime"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/gtfs"
)

func NewDatabaseFromFilesystem(path string, t time.Time) (*Database, *model.Base) {
	tz, _ := time.LoadLocation("America/Montreal")

	raw, err := gtfs.RawFilesystem(path)
	if err != nil {
		panic(err)
	}

	base, err := gtfs.NewBase(raw, gtfs.BaseOptions{
		Time:       t,
		TimeZone:   tz,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	})

	if err != nil {
		panic(err)
	}

	database := NewDatabase(base, tz)
	runtime.GC()

	return database, base
}
