package daos

import (
	"errors"
	"luvsic3/uvid/models"
	"luvsic3/uvid/tools"
	"time"

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

	err = db.AutoMigrate(&models.User{}, &models.Session{}, &models.Event{}, &models.PerformanceSpan{}, &models.HTTPSpan{}, &models.JSError{})
	if err != nil {
		panic(err)
	}

	return &Dao{
		db,
	}
}

func (dao *Dao) TimeRange(start time.Time, end time.Time) *gorm.DB {
	return dao.DB.Where("created_at >= ? AND created_at < ?", start, end).Session(&gorm.Session{})
}

const rootUserName = "root"

func (dao *Dao) InitializeDB() error {
	rootUser := models.User{
		Name: rootUserName,
	}
	err := dao.DB.First(&rootUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pass, err := tools.HashPassword(rootUserName)
		if err != nil {
			return err
		}
		err = dao.DB.Create(&models.User{
			Name:     rootUserName,
			Password: string(pass),
		}).Error
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
