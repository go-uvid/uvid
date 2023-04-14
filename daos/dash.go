package daos

import (
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/models"
	"luvsic3/uvid/tools"

	"gorm.io/gorm"
)

func IsPageView(db *gorm.DB) *gorm.DB {
	return db.Model(models.PageView{})
}

func (dao *Dao) FindPageViewCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Scopes(IsPageView).
		Count(&count).Error
	return count, err
}

type IntervalData = dtos.IntervalData

const hourAndCountColumn = "strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'localtime')) as x, COUNT(*) as y"
const dayAndCountColumn = "strftime('%Y-%m-%d', datetime(created_at, 'localtime')) as x, COUNT(*) as y"
const hourAndUniqueCountColumn = "strftime('%Y-%m-%d %H:00:00', datetime(created_at, 'localtime')) as x, COUNT(DISTINCT session_uuid) as y"
const dayAndUniqueCountColumn = "strftime('%Y-%m-%d', datetime(created_at, 'localtime')) as x, COUNT(DISTINCT session_uuid) as y"

// findPageViews returns the number of page views in the given time range
func (dao *Dao) FindPageViewInterval(db *gorm.DB, byHour bool) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Scopes(IsPageView).
		Select(tools.Ternary(byHour, hourAndCountColumn, dayAndCountColumn)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error

	return results, err
}

// findUniqueVisitors returns the number of unique visitors in the given time range
func (dao *Dao) FindUniqueVisitorInterval(db *gorm.DB, byHour bool) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Scopes(IsPageView).
		Select(tools.Ternary(byHour, hourAndUniqueCountColumn, dayAndUniqueCountColumn)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error
	return results, err
}

func (dao *Dao) FindUniqueVisitorCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Scopes(IsPageView).
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
		Select("name as x, COUNT(*) as y").
		Group("name").
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
func (dao *Dao) FindJSErrorInterval(db *gorm.DB, byHour bool) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Model(&models.JSError{}).
		Select(tools.Ternary(byHour, hourAndCountColumn, dayAndCountColumn)).
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
func (dao *Dao) FindHTTPErrorInterval(db *gorm.DB, byHour bool) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Model(&models.HTTP{}).
		Select(tools.Ternary(byHour, hourAndCountColumn, dayAndCountColumn)).
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
