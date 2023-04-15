package daos

import (
	"fmt"
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/models"
	"luvsic3/uvid/tools"

	"gorm.io/gorm"
)

type TimeFormat string

const (
	FormatHour  TimeFormat = "%Y-%m-%d %H:00:00"
	FormatDay   TimeFormat = "%Y-%m-%d"
	FormatMonth TimeFormat = "%Y-%m"
	FormatYear  TimeFormat = "%Y"
)

func unitToTimeFormat(unit tools.Unit) TimeFormat {
	switch unit {
	case tools.UnitHour:
		return FormatHour
	case tools.UnitDay:
		return FormatDay
	case tools.UnitMonth:
		return FormatMonth
	case tools.UnitYear:
		return FormatYear
	default:
		return FormatDay
	}
}

const DistinctSession = "DISTINCT session_uuid"

func isPageView(db *gorm.DB) *gorm.DB {
	return db.Model(models.PageView{})
}

func selectColumn(uniqueSession bool, unit tools.Unit) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fmt.Sprintf("strftime('%s', datetime(created_at, 'localtime')) as x, COUNT(%s) as y", unitToTimeFormat(unit), tools.Ternary(uniqueSession, DistinctSession, "*")))
	}
}

func (dao *Dao) FindPageViewCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Scopes(isPageView).
		Count(&count).Error
	return count, err
}

type IntervalData = dtos.IntervalData

// findPageViews returns the number of page views in the given time range
func (dao *Dao) FindPageViewInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Scopes(isPageView, selectColumn(false, unit)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error

	return results, err
}

// findUniqueVisitors returns the number of unique visitors in the given time range
func (dao *Dao) FindUniqueVisitorInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Scopes(isPageView, selectColumn(true, unit)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error
	return results, err
}

func (dao *Dao) FindUniqueVisitorCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Scopes(isPageView).
		Distinct("session_uuid").
		Count(&count).Error
	return count, err
}

// findAveragePerformance returns the average performance spans in the given time range
func (dao *Dao) FindAveragePerformanceInterval(db *gorm.DB) ([]IntervalData, error) {
	var results []IntervalData
	err := db.Model(&models.Performance{}).
		Select("name as x, AVG(value) as y").
		Group("name").
		Scan(&results).Error
	return results, err
}

// findEvents returns the number of events in the given time range
func (dao *Dao) FindEventInterval(db *gorm.DB) ([]IntervalData, error) {
	var results []IntervalData
	err := db.Model(&models.Event{}).
		Select("action as x, COUNT(*) as y").
		Group("action").
		Order("y DESC").
		Scan(&results).Error
	return results, err
}

func (dao *Dao) FindJSErrorCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.JSError{}).
		Count(&count).Error
	return count, err
}

// findJSErrors returns the number of JS errors in the given time range
func (dao *Dao) FindJSErrorInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Model(&models.JSError{}).
		Scopes(selectColumn(false, unit)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error
	return results, err
}

func (dao *Dao) FindHTTPErrorCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(models.HTTP{}).Where("status < ? or status > ?", 200, 299).Count(&count).Error
	return count, err
}

// findHTTPErrors returns the number of HTTP errors in the given time range
func (dao *Dao) FindHTTPErrorInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Model(&models.HTTP{}).
		Scopes(selectColumn(false, unit)).
		Where("status < ? or status > ?", 200, 299).
		Group("x").
		Order("x ASC").
		Scan(&results).Error
	return results, err
}

// findHTTPErrors returns the number of HTTP errors in the given time range
func (dao *Dao) GetUserByName(name string) (models.User, error) {
	user := models.User{
		Name: name,
	}

	if err := dao.DB.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// findHTTPErrors returns the number of HTTP errors in the given time range
func (dao *Dao) UpdateUserPassword(name string, password string) error {
	user := models.User{
		Name: name,
	}

	return dao.DB.Where(user).Update("password", password).Error
}
