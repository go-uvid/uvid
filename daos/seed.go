package daos

import (
	"fmt"
	"luvsic3/uvid/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(dsn string) {
	// Connect to the SQLite database
	dao := New(dsn)
	db := dao.DB

	var count int64
	db.Model(&models.Session{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already seeded. Skipping...")
		return
	}
	fmt.Println("Start seeding database...")
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
			Referrer:   getRandomDomain(),
			Meta:       "{}",
			Model:      gorm.Model{CreatedAt: getRandomTime(7)},
			UUID:       uuid.New(),
		}
		if err := db.Create(&session).Error; err != nil {
			panic(err)
		}

		for j := 0; j < size; j++ {
			perfMetric := models.PageView{
				URL:         getRandomPath(),
				SessionUUID: session.UUID,
				Session:     session,
				Model:       gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&perfMetric)
		}

		for j := 0; j < size; j++ {
			perfMetric := models.Performance{
				Name:        getRandomPerfName(),
				Value:       getRandomPerfValue(),
				URL:         getRandomPath(),
				SessionUUID: session.UUID,
				Session:     session,
				Model:       gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&perfMetric)
		}

		// Generate 200 JSError records for this session
		for j := 0; j < size; j++ {
			error := &models.JSError{
				Name:        "error",
				Message:     "Something went wrong",
				Stack:       "Error: Something went wrong",
				SessionUUID: session.UUID,
				Session:     session,
				Model:       gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&error)
		}

		// Generate 200 HTTPMetric records for this session
		for j := 0; j < size; j++ {
			metric := &models.HTTP{
				Resource:    getRandomDomain(),
				Method:      randomHttpMethod(),
				Headers:     "Content-Type: text/html",
				Status:      rand.Intn(400) + 100,
				SessionUUID: session.UUID,
				Session:     session,
				Model:       gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&metric)
		}

		// Generate 200 Event records for this session
		for j := 0; j < size; j++ {
			event := &models.Event{
				Action:      randomEventAction(),
				Value:       "",
				SessionUUID: session.UUID,
				Session:     session,
				Model:       gorm.Model{CreatedAt: getRandomTime(7)},
			}
			db.Create(&event)
		}
	}
	fmt.Print("Database seeded successfully!")
}

// getRandomPerfName generates a random name for a Performance
func getRandomPerfName() string {
	names := []string{models.LCP, models.CLS, models.FID}
	return names[rand.Intn(len(names))]
}

// getRandomPerfValue generates a random value for a Performance
func getRandomPerfValue() float64 {
	return rand.Float64() * 10
}

// getRandomDomain generates a random URL for a Performance
func getRandomDomain() string {
	urls := []string{"https://example.com", "https://google.com", "https://github.com", "https://stackoverflow.com", "https://wikipedia.org"}
	return urls[rand.Intn(len(urls))]
}

// getRandomPath gererates a random url path
func getRandomPath() string {
	paths := []string{"/", "/about", "/contact", "/login", "/register", "/dashboard", "/profile", "/settings", "/admin", "/admin/dashboard", "/admin/users"}
	return "https://example.com" + paths[rand.Intn(len(paths))]
}

func getRandomTime(daysOffset int) time.Time {
	hour := -rand.Intn(24)
	return time.Now().AddDate(0, 0, -rand.Intn(daysOffset)).Add(time.Duration(hour) * time.Hour)
}

func randomHttpMethod() string {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions, http.MethodConnect, http.MethodTrace}
	return methods[rand.Intn(len(methods))]
}

func randomEventAction() string {
	names := []string{"register", "login", "logout", "click", "view", "add", "remove", "update"}
	return names[rand.Intn(len(names))]
}
