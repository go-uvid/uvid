package daos

import (
	"luvsic3/uvid/models"
	"luvsic3/uvid/tools"

	"gorm.io/gorm"
)

func (dao *Dao) FindPageViewCount(db *gorm.DB) int64 {
	var count int64
	db.Model(models.PerformanceMetric{}).
		Where("name = ?", models.LCP).
		Count(&count)
	return count
}

type IntervalData struct {
	X string
	Y int64
}

const hourAndCountColumn = "strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'localtime')) as x, COUNT(*) as y"
const dayAndCountColumn = "strftime('%Y-%m-%d', datetime(created_at, 'localtime')) as x, COUNT(*) as y"
const hourAndUniqueCountColumn = "strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'localtime')) as x, COUNT(DISTINCT session_id) as y"
const dayAndUniqueCountColumn = "strftime('%Y-%m-%d', datetime(created_at, 'localtime')) as x, COUNT(DISTINCT session_id) as y"

// findPageViews returns the number of page views in the given time range
func (dao *Dao) FindPageViewInterval(db *gorm.DB, byHour bool) []IntervalData {
	var results []IntervalData

	// Build the query to group by hour, count distinct sessions with LCP metric
	db.Model(&models.PerformanceMetric{}).
		Select(tools.Ternary(byHour, hourAndCountColumn, dayAndCountColumn)).
		Where("name = ?", models.LCP).
		Group("x").
		Order("x ASC").
		Scan(&results)
	return results
}

// findUniqueVisitors returns the number of unique visitors in the given time range
func (dao *Dao) FindUniqueVisitorInterval(db *gorm.DB, byHour bool) []IntervalData {
	var results []IntervalData

	// Build the query to group by hour, count distinct sessions with LCP metric
	db.Model(&models.PerformanceMetric{}).
		Select(tools.Ternary(byHour, hourAndUniqueCountColumn, dayAndUniqueCountColumn)).
		Where("name = ?", models.LCP).
		Group("x").
		Order("x ASC").
		Scan(&results)
	return results
}

func (dao *Dao) FindUniqueVisitorCount(db *gorm.DB) int64 {
	var count int64
	db.Model(models.PerformanceMetric{}).
		Distinct("session_id").
		Where("name = ?", models.LCP).
		Count(&count)
	return count
}

// findAveragePerformance returns the average performance metrics in the given time range
func (dao *Dao) FindAveragePerformanceInterval(db *gorm.DB) []IntervalData {
	var results []IntervalData
	db.Model(&models.PerformanceMetric{}).
		Select("name as x, AVG(value) as y").
		Group("name").
		Scan(&results)
	return results
}

// findEvents returns the number of events in the given time range
func (dao *Dao) FindEventInterval(db *gorm.DB) []IntervalData {
	var results []IntervalData
	db.Model(&models.Event{}).
		Select("name as x, COUNT(*) as y").
		Group("name").
		Scan(&results)
	return results
}

func (dao *Dao) FindJSErrorCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&models.JSError{}).
		Count(&count)
	return count
}

// findJSErrors returns the number of JS errors in the given time range
func (dao *Dao) FindJSErrorInterval(db *gorm.DB, byHour bool) []IntervalData {
	var results []IntervalData

	db.Model(&models.JSError{}).
		Select(tools.Ternary(byHour, hourAndCountColumn, dayAndCountColumn)).
		Group("x").
		Order("x ASC").
		Scan(&results)
	return results
}

func (dao *Dao) FindHTTPErrorCount(db *gorm.DB) int64 {
	var count int64
	db.Model(models.HTTPMetric{}).Where("status < ? or status > ?", 200, 299).Count(&count)
	return count
}

// findHTTPErrors returns the number of HTTP errors in the given time range
func (dao *Dao) FindHTTPErrorInterval(db *gorm.DB, byHour bool) []IntervalData {
	var results []IntervalData

	db.Model(&models.HTTPMetric{}).
		Select(tools.Ternary(byHour, hourAndCountColumn, dayAndCountColumn)).
		Where("status < ? or status > ?", 200, 299).
		Group("x").
		Order("x ASC").
		Scan(&results)
	return results
}