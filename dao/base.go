package dao

import (
	"luvsic3/uvid/model"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func New() *Dao {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{}, &model.Session{}, &model.Event{}, &model.PerformanceMetric{}, &model.HTTPMetric{}, &model.JSError{})

	return &Dao{
		db,
	}
}

func (dao *Dao) TimeRange(start time.Time, end time.Time) *gorm.DB {
	return dao.db.Where("created_at >= ? AND created_at < ?", start, end)
}

func (dao *Dao) FindUniqueVisitorCount() int64 {
	var count int64
	dao.db.Model(model.PerformanceMetric{}).
		Select("COUNT(DISTINCT(session_id))").
		Where("name = ?", "LCP").
		Count(&count)
	return count
}
