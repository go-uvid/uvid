package daos

import (
	"fmt"
	"math/rand"
	"net/http"
	"rick-you/uvid/models"
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
			UA:         randomUserAgent(),
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
	urls := []string{"https://google.com", "https://github.com", "https://stackoverflow.com", "https://wikipedia.org"}
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

// random user agent, chrome firefox safari etc
func randomUserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
	}
	return agents[rand.Intn(len(agents))]
}
