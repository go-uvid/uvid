package daos_test

import (
	"luvsic3/uvid/daos"
	"luvsic3/uvid/models"
	"luvsic3/uvid/tools"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFindPageViewInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	results, _ := dao.FindPageViewInterval(dao.SpanFilter(time.Time{}, time.Now()), tools.UnitHour)
	assert.Empty(t, results)

	startTime := truncateToday()
	endTime := startTime.AddDate(0, 0, 1)

	db.Create(&models.PageView{URL: "", Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PageView{URL: "", Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PageView{URL: "", Model: gorm.Model{CreatedAt: endTime}})

	results, _ = dao.FindPageViewInterval(dao.SpanFilter(startTime, startTime.Add(time.Hour)), tools.UnitHour)
	result := results[0]

	assert.Len(t, results, 1)
	assert.Equal(t, result.X, startTime.Truncate(time.Hour).Format(time.DateTime))
	assert.Equal(t, result.Y, int64(2))
}

func TestFindAveragePerformanceInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different values
	startTime := truncateToday()
	endTime := startTime.AddDate(0, 0, 1)
	db.Create(&models.Performance{Name: models.LCP, Value: 1.2, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.Performance{Name: models.LCP, Value: 2.1, Model: gorm.Model{CreatedAt: endTime}})
	db.Create(&models.Performance{Name: models.LCP, Value: 3.2, Model: gorm.Model{CreatedAt: endTime}})

	// Test by finding the average performance interval for the last hour
	results, _ := dao.FindAveragePerformanceInterval(dao.SpanFilter(endTime, endTime.Add(time.Hour)))
	result := results[0]

	// Check that the result is as expected
	assert.Len(t, results, 1)
	assert.Equal(t, result.X, models.LCP)
	assert.InDelta(t, result.Y, float64((200+300)/2), 0.01)
}

func TestFindHTTPErrorInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some HTTP span with different statuses
	startTime := truncateToday()
	db.Create(&models.HTTP{Status: 200, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.HTTP{Status: 204, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.HTTP{Status: 300, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour + time.Minute)}})
	db.Create(&models.HTTP{Status: 400, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour + time.Minute)}})
	db.Create(&models.HTTP{Status: 500, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour * 2)}})

	// Test by finding the HTTP error interval for the last hour
	results, _ := dao.FindHTTPErrorInterval(dao.SpanFilter(startTime.Add(time.Hour), startTime.Add(time.Hour*2)), tools.UnitHour)
	result := results[0]

	// Check that the result is as expected
	assert.Len(t, results, 1)
	assert.Equal(t, result.X, startTime.Add(time.Hour).Truncate(time.Hour).Format(time.DateTime))
	assert.Equal(t, result.Y, int64(2))
}

func TestFindUniqueVisitorInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()
	db.Create(&models.PageView{URL: "", SessionUUID: uuid1, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PageView{URL: "", SessionUUID: uuid1, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.PageView{URL: "", SessionUUID: uuid2, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.PageView{URL: "", SessionUUID: uuid3, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.PageView{URL: "", SessionUUID: uuid3, Model: gorm.Model{CreatedAt: endTime}})

	// Test by finding the unique visitor interval for the last 2 hours
	tr := dao.SpanFilter(startTime, endTime)
	results, _ := dao.FindUniqueVisitorInterval(tr, tools.UnitDay)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, startTime.Format(time.DateOnly))
	assert.Equal(t, results[0].Y, int64(1))

	assert.Equal(t, results[1].X, midTime.Format(time.DateOnly))
	assert.Equal(t, results[1].Y, int64(2))
}

// test dao.FindEventInterval
func TestFindEventInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	uuid1 := uuid.New()
	const registerAction = "register"
	const loginAction = "login"

	db.Create(&models.Event{Action: registerAction, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.Event{Action: registerAction, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.Event{Action: registerAction, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.Event{Action: loginAction, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.Event{Action: loginAction, SessionUUID: uuid1, Model: gorm.Model{CreatedAt: endTime}})

	tr := dao.SpanFilter(startTime, endTime)
	results, _ := dao.FindEventInterval(tr)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, registerAction)
	assert.Equal(t, results[0].Y, int64(3))

	assert.Equal(t, results[1].X, loginAction)
	assert.Equal(t, results[1].Y, int64(1))
}

// test dao.FindJSErrorInterval
func TestFindJSErrorInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	uuid1 := uuid.New()
	db.Create(&models.JSError{
		Name:        "error",
		Message:     "Something went wrong",
		Stack:       "Error: Something went wrong",
		SessionUUID: uuid1,
		Model:       gorm.Model{CreatedAt: startTime},
	})
	db.Create(&models.JSError{
		Name:        "error",
		Message:     "Something went wrong",
		Stack:       "Error: Something went wrong",
		SessionUUID: uuid1,
		Model:       gorm.Model{CreatedAt: startTime.Add(time.Hour)},
	})
	db.Create(&models.JSError{
		Name:        "error",
		Message:     "Something went wrong",
		Stack:       "Error: Something went wrong",
		SessionUUID: uuid1,
		Model:       gorm.Model{CreatedAt: midTime},
	})
	db.Create(&models.JSError{
		Name:        "error",
		Message:     "Something went wrong",
		Stack:       "Error: Something went wrong",
		SessionUUID: uuid1,
		Model:       gorm.Model{CreatedAt: midTime},
	})
	db.Create(&models.JSError{
		Name:        "error",
		Message:     "Something went wrong",
		Stack:       "Error: Something went wrong",
		SessionUUID: uuid1,
		Model:       gorm.Model{CreatedAt: endTime},
	})

	// Test by finding the unique visitor interval for the last 2 hours
	tr := dao.SpanFilter(startTime, endTime)
	results, _ := dao.FindJSErrorInterval(tr, tools.UnitDay)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, startTime.Format(time.DateOnly))
	assert.Equal(t, results[0].Y, int64(2))

	assert.Equal(t, results[1].X, midTime.Format(time.DateOnly))
	assert.Equal(t, results[1].Y, int64(2))
}

func truncateToday() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

func NewInMemoryDao() *daos.Dao {
	return daos.New(":memory:")
}
