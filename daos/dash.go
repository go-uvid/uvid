package daos

import (
	"fmt"
	"rick-you/uvid/dtos"
	"rick-you/uvid/models"
	"rick-you/uvid/tools"

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

func (dao *Dao) FindPageViews(db *gorm.DB) ([]dtos.PageViewDTO, error) {
	var results []dtos.PageViewDTO

	err := db.Scopes(isPageView).Scan(&results).Error

	return results, err
}

func (dao *Dao) FindPageViewInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Scopes(isPageView, selectColumn(false, unit)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error

	return results, err
}

func (dao *Dao) FindSessions(db *gorm.DB) ([]dtos.SessionDTO, error) {
	var results []dtos.SessionDTO

	err := db.Model(&models.Session{}).Scan(&results).Error
	return results, err
}

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

func (dao *Dao) FindAveragePerformanceInterval(db *gorm.DB) ([]IntervalData, error) {
	var results []IntervalData
	err := db.Model(&models.Performance{}).
		Select("name as x, AVG(value) as y").
		Group("name").
		Scan(&results).Error
	return results, err
}

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

func (dao *Dao) FindJSErrors(db *gorm.DB) ([]dtos.ErrorDTO, error) {
	var results []dtos.ErrorDTO

	err := db.Model(&models.JSError{}).Scan(&results).Error
	return results, err
}

func (dao *Dao) FindJSErrorInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Model(&models.JSError{}).
		Scopes(selectColumn(false, unit)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error
	return results, err
}

func isHTTPError(db *gorm.DB) *gorm.DB {
	return db.Model(models.HTTP{}).Where("status < ? or status > ?", 200, 299)
}

func (dao *Dao) FindHTTPErrorCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Scopes(isHTTPError).Count(&count).Error
	return count, err
}

func (dao *Dao) FindHTTPErrors(db *gorm.DB) ([]dtos.HTTPDTO, error) {
	var results []dtos.HTTPDTO

	err := db.Scopes(isHTTPError).Scan(&results).Error
	return results, err
}

func (dao *Dao) FindHTTPErrorInterval(db *gorm.DB, unit tools.Unit) ([]IntervalData, error) {
	var results []IntervalData

	err := db.Scopes(isHTTPError, selectColumn(false, unit)).
		Group("x").
		Order("x ASC").
		Scan(&results).Error
	return results, err
}

func (dao *Dao) GetUserByName(name string) (models.User, error) {
	user := models.User{
		Name: name,
	}

	if err := dao.DB.First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (dao *Dao) ChangeUserPassword(name string, currentPassword, newPassword string) error {
	user := models.User{
		Name: name,
	}
	dao.DB.Where("name = ?", user.Name).First(&user)
	if err := tools.ComparePassword(user.Password, currentPassword); err != nil {
		return err
	}
	hashPass, err := tools.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return dao.DB.Model(&models.User{}).Where(user).Update("password", string(hashPass)).Error
}
