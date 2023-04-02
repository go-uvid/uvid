package daos_test

import (
	"luvsic3/uvid/daos"
	"luvsic3/uvid/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFindPageViewInterval(t *testing.T) {
	dao := daos.NewInMemoryDao()
	db := dao.DB()

	results := dao.FindPageViewInterval(dao.TimeRange(time.Time{}, time.Now()), true)
	assert.Empty(t, results)

	startTime := truncateToday()
	endTime := startTime.AddDate(0, 0, 1)

	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 100, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 200, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 300, Model: gorm.Model{CreatedAt: endTime}})

	results = dao.FindPageViewInterval(dao.TimeRange(startTime, startTime.Add(time.Hour)), true)
	result := results[0]

	assert.Len(t, results, 1)
	assert.Equal(t, result.X, startTime.Truncate(time.Hour).Format(time.DateTime))
	assert.Equal(t, result.Y, int64(2))
}

func TestFindAveragePerformanceInterval(t *testing.T) {
	dao := daos.NewInMemoryDao()
	db := dao.DB()

	// Create some performance span with different values
	startTime := truncateToday()
	endTime := startTime.AddDate(0, 0, 1)
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 100, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 200, Model: gorm.Model{CreatedAt: endTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 300, Model: gorm.Model{CreatedAt: endTime}})

	// Test by finding the average performance interval for the last hour
	results := dao.FindAveragePerformanceInterval(dao.TimeRange(endTime, endTime.Add(time.Hour)))
	result := results[0]

	// Check that the result is as expected
	assert.Len(t, results, 1)
	assert.Equal(t, result.X, models.LCP)
	assert.InDelta(t, result.Y, float64((200+300)/2), 0.01)
}

func TestFindHTTPErrorInterval(t *testing.T) {
	dao := daos.NewInMemoryDao()
	db := dao.DB()

	// Create some HTTP span with different statuses
	startTime := truncateToday()
	db.Create(&models.HTTPSpan{Status: 200, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.HTTPSpan{Status: 204, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.HTTPSpan{Status: 300, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour + time.Minute)}})
	db.Create(&models.HTTPSpan{Status: 400, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour + time.Minute)}})
	db.Create(&models.HTTPSpan{Status: 500, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour * 2)}})

	// Test by finding the HTTP error interval for the last hour
	results := dao.FindHTTPErrorInterval(dao.TimeRange(startTime.Add(time.Hour), startTime.Add(time.Hour*2)), true)
	result := results[0]

	// Check that the result is as expected
	assert.Len(t, results, 1)
	assert.Equal(t, result.X, startTime.Add(time.Hour).Truncate(time.Hour).Format(time.DateTime))
	assert.Equal(t, result.Y, int64(2))
}

func TestFindUniqueVisitorInterval(t *testing.T) {
	dao := daos.NewInMemoryDao()
	db := dao.DB()

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 100, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 300, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 200, SessionUUID: uuid2, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 200, SessionUUID: uuid3, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.PerformanceSpan{Name: models.LCP, Value: 400, SessionUUID: uuid3, Model: gorm.Model{CreatedAt: endTime}})

	// Test by finding the unique visitor interval for the last 2 hours
	tr := dao.TimeRange(startTime, endTime)
	results := dao.FindUniqueVisitorInterval(tr, false)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, startTime.Format(time.DateOnly))
	assert.Equal(t, results[0].Y, int64(1))

	assert.Equal(t, results[1].X, midTime.Format(time.DateOnly))
	assert.Equal(t, results[1].Y, int64(2))
}

func truncateToday() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}
