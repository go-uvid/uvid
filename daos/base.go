package daos

import (
	"errors"
	"time"

	"github.com/rick-you/uvid/models"
	"github.com/rick-you/uvid/tools"

	"github.com/google/uuid"
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

	err = db.AutoMigrate(&models.User{}, &models.Session{}, &models.Event{}, &models.Performance{}, &models.HTTP{}, &models.JSError{}, &models.PageView{}, &models.Config{})
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

const (
	ADMIN_NAME = "root"
	ADMIN_PASS = "uvid"
)

func (dao *Dao) InitializeDB() error {
	rootUser := models.User{
		Name: ADMIN_NAME,
	}
	err := dao.DB.First(&rootUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pass, err := tools.HashPassword(ADMIN_PASS)
		if err != nil {
			return err
		}
		dao.DB.Transaction(func(tx *gorm.DB) error {
			if err = tx.Create(&models.User{
				Name:     ADMIN_NAME,
				Password: string(pass),
			}).Error; err != nil {
				return err
			}

			if err = tx.Create(&models.Config{
				Key:   "jwt_secret",
				Value: uuid.NewString(),
			}).Error; err != nil {
				return err
			}

			return nil
		})

		return nil
	}
	return err
}

func (dao *Dao) GetAllConfigs() (map[string]string, error) {
	var configs []models.Config

	// Query for all rows in the config table
	if err := dao.DB.Find(&configs).Error; err != nil {
		return nil, err
	}

	// Transform the slice of Config structs into a key-value map
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	return configMap, nil
}
