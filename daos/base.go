package daos

import (
	"luvsic3/uvid/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func New(dsn string) *Dao {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Session{}, &models.Event{}, &models.PerformanceSpan{}, &models.HTTPSpan{}, &models.JSError{})
	if err != nil {
		panic(err)
	}

	return &Dao{
		db,
	}
}

func NewInMemoryDao() *Dao {
	return New(":memory:")
}

func (dao *Dao) DB() *gorm.DB {
	return dao.db
}

func (dao *Dao) TimeRange(start time.Time, end time.Time) *gorm.DB {
	return dao.db.Where("created_at >= ? AND created_at < ?", start, end).Session(&gorm.Session{})
}
