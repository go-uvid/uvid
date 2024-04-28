package daos_test

import (
	"testing"
	"time"

	"github.com/go-uvid/uvid/daos"
	"github.com/go-uvid/uvid/models"
	"github.com/go-uvid/uvid/tools"

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
	assert.Equal(t, result.Y, float64(2))
}

func TestFindAveragePerformanceInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different values
	startTime := truncateToday()
	endTime := startTime.AddDate(0, 0, 1)
	db.Create(&models.Performance{Name: models.LCP, Value: 1.1, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.Performance{Name: models.LCP, Value: 2.2, Model: gorm.Model{CreatedAt: endTime}})
	db.Create(&models.Performance{Name: models.LCP, Value: 3.3, Model: gorm.Model{CreatedAt: endTime}})

	// Test by finding the average performance interval for the last hour
	results, _ := dao.FindAveragePerformanceInterval(dao.SpanFilter(endTime, endTime.Add(time.Hour)))
	result := results[0]

	// Check that the result is as expected
	assert.Len(t, results, 1)
	assert.Equal(t, result.X, models.LCP)
	assert.InDelta(t, result.Y, float64((2.2+3.3)/2), 0.01)
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
	assert.Equal(t, result.Y, float64(2))
}

func TestFindUniqueVisitorInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	id1 := uint(1)
	id2 := uint(2)
	id3 := uint(3)
	db.Create(&models.PageView{URL: "", SessionID: id1, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.PageView{URL: "", SessionID: id1, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.PageView{URL: "", SessionID: id2, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.PageView{URL: "", SessionID: id3, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.PageView{URL: "", SessionID: id3, Model: gorm.Model{CreatedAt: endTime}})

	// Test by finding the unique visitor interval for the last 2 hours
	tr := dao.SpanFilter(startTime, endTime)
	results, _ := dao.FindUniqueVisitorInterval(tr, tools.UnitDay)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, startTime.Format(time.DateOnly))
	assert.Equal(t, results[0].Y, float64(1))

	assert.Equal(t, results[1].X, midTime.Format(time.DateOnly))
	assert.Equal(t, results[1].Y, float64(2))
}

// test dao.FindEventInterval
func TestFindEventInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	id1 := uint(1)
	const registerAction = "register"
	const loginAction = "login"

	db.Create(&models.Event{Action: registerAction, SessionID: id1, Model: gorm.Model{CreatedAt: startTime}})
	db.Create(&models.Event{Action: registerAction, SessionID: id1, Model: gorm.Model{CreatedAt: startTime.Add(time.Hour)}})
	db.Create(&models.Event{Action: registerAction, SessionID: id1, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.Event{Action: loginAction, SessionID: id1, Model: gorm.Model{CreatedAt: midTime}})
	db.Create(&models.Event{Action: loginAction, SessionID: id1, Model: gorm.Model{CreatedAt: endTime}})

	tr := dao.SpanFilter(startTime, endTime)
	results, _ := dao.FindEventInterval(tr)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, registerAction)
	assert.Equal(t, results[0].Y, float64(3))

	assert.Equal(t, results[1].X, loginAction)
	assert.Equal(t, results[1].Y, float64(1))
}

// test dao.FindJSErrorInterval
func TestFindJSErrorInterval(t *testing.T) {
	dao := NewInMemoryDao()
	db := dao.DB

	// Create some performance span with different session ids
	startTime := truncateToday()
	midTime := startTime.AddDate(0, 0, 1)
	endTime := startTime.AddDate(0, 0, 3)
	id1 := uint(1)
	db.Create(&models.JSError{
		Name:      "error",
		Message:   "Something went wrong",
		Stack:     "Error: Something went wrong",
		SessionID: id1,
		Model:     gorm.Model{CreatedAt: startTime},
	})
	db.Create(&models.JSError{
		Name:      "error",
		Message:   "Something went wrong",
		Stack:     "Error: Something went wrong",
		SessionID: id1,
		Model:     gorm.Model{CreatedAt: startTime.Add(time.Hour)},
	})
	db.Create(&models.JSError{
		Name:      "error",
		Message:   "Something went wrong",
		Stack:     "Error: Something went wrong",
		SessionID: id1,
		Model:     gorm.Model{CreatedAt: midTime},
	})
	db.Create(&models.JSError{
		Name:      "error",
		Message:   "Something went wrong",
		Stack:     "Error: Something went wrong",
		SessionID: id1,
		Model:     gorm.Model{CreatedAt: midTime},
	})
	db.Create(&models.JSError{
		Name:      "error",
		Message:   "Something went wrong",
		Stack:     "Error: Something went wrong",
		SessionID: id1,
		Model:     gorm.Model{CreatedAt: endTime},
	})

	// Test by finding the unique visitor interval for the last 2 hours
	tr := dao.SpanFilter(startTime, endTime)
	results, _ := dao.FindJSErrorInterval(tr, tools.UnitDay)

	assert.Len(t, results, 2)
	// Check that the result is as expected
	assert.Equal(t, results[0].X, startTime.Format(time.DateOnly))
	assert.Equal(t, results[0].Y, float64(2))

	assert.Equal(t, results[1].X, midTime.Format(time.DateOnly))
	assert.Equal(t, results[1].Y, float64(2))
}

func truncateToday() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

func NewInMemoryDao() *daos.Dao {
	return daos.New(":memory:")
}
