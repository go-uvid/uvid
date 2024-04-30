package daos

import (
	"time"

	"github.com/rick-you/uvid/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dao struct {
	DB *gorm.DB
}

func New(dsn string) *Dao {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Session{}, &models.Event{}, &models.Performance{}, &models.HTTP{}, &models.JSError{}, &models.PageView{})
	if err != nil {
		panic(err)
	}

	return &Dao{
		db,
	}
}

func (dao *Dao) SpanFilter(start time.Time, end time.Time) *gorm.DB {
	// FIXME looks like: sqlite store time in local time, but client request in ISO time, so we need to convert it to local time before query
	return dao.DB.Where("created_at >= ? AND created_at < ?", start.Local(), end.Local()).Session(&gorm.Session{})
}
