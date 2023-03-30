package tools

import (
	"fmt"
	"luvsic3/uvid/daos"
	"luvsic3/uvid/models"
	"math/rand"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func Seed(dsn string) {
	// Connect to the SQLite database
	dao := daos.New(dsn)
	db := dao.DB()

	var count int64
	db.Model(&models.Session{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already seeded. Skipping...")
		return
	}
	const size = 100
	// Seed the database with 10 Session records
	for i := 0; i < 10; i++ {
		session := models.Session{
			UA:         "Mozilla/5.0 (Windows NT 10.0; Win64)",
			Language:   "en-US",
			IP:         "192.168.0.1",
			AppVersion: "1.0.0",
			URL:        "https://example.com",
			Screen:     "1920x1080",
			Referrer:   getRandomURL(),
			Meta:       "{}",
			Model:      gorm.Model{CreatedAt: getRandomTime(7)},
		}
		if err := db.Create(&session).Error; err != nil {
			panic(err)
		}

		for j := 0; j < size; j++ {
			perfMetric := models.PerformanceMetric{
				Name:      getRandomPerfName(),
				Value:     getRandomPerfValue(),
				URL:       getRandomURL(),
				SessionID: session.ID,
				Model:     gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&perfMetric)
		}

		// Generate 200 JSError records for this session
		for j := 0; j < size; j++ {
			error := &models.JSError{
				Error:   fmt.Sprintf("Error %d", j),
				Session: session,
				Model:   gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&error)
		}

		// Generate 200 HTTPMetric records for this session
		for j := 0; j < size; j++ {
			metric := &models.HTTPMetric{
				URL:      session.URL + "page",
				Method:   randomHttpMethod(),
				Headers:  "Content-Type: text/html",
				Status:   rand.Intn(400) + 100,
				Response: "",
				Session:  session,
				Model:    gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&metric)
		}

		// Generate 200 Event records for this session
		for j := 0; j < size; j++ {
			event := &models.Event{
				Name:    randomEventName(),
				Value:   "",
				Session: session,
				Model:   gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&event)
		}
	}
}

// getRandomPerfName generates a random name for a PerformanceMetric
func getRandomPerfName() string {
	names := []string{models.LCP, models.CLS, models.FID}
	return names[rand.Intn(len(names))]
}

// getRandomPerfValue generates a random value for a PerformanceMetric
func getRandomPerfValue() float64 {
	return rand.Float64() * 10
}

// getRandomURL generates a random URL for a PerformanceMetric
func getRandomURL() string {
	urls := []string{"https://example.com", "https://google.com", "https://github.com", "https://stackoverflow.com", "https://wikipedia.org"}
	return urls[rand.Intn(len(urls))]
}

func getRandomTime(daysOffset int) time.Time {
	return time.Now().AddDate(0, 0, -rand.Intn(daysOffset))
}

func randomHttpMethod() string {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect, http.MethodTrace}
	return methods[rand.Intn(len(methods))]
}

func randomEventName() string {
	names := []string{"register", "login", "logout", "click", "view", "add", "remove", "update"}
	return names[rand.Intn(len(names))]
}
