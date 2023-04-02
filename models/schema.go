package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"uniqueIndex;type:byte"`
	UA         string    `gorm:"not null"`
	Language   string    `gorm:"not null"`
	IP         string    `gorm:"not null"`
	AppVersion string
	URL        string `gorm:"not null"`
	Screen     string `gorm:"not null"`
	Referrer   string
	Meta       string

	JSErrors        []JSError         `gorm:"foreignKey:SessionUUID"`
	HTTPSpan        []HTTPSpan        `gorm:"foreignKey:SessionUUID"`
	Events          []Event           `gorm:"foreignKey:SessionUUID"`
	PerformanceSpan []PerformanceSpan `gorm:"foreignKey:SessionUUID"`
}

type JSError struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Message string `gorm:"not null"`
	Stack   string `gorm:"not null"`
	Cause   string

	SessionUUID uuid.UUID `gorm:"not null"`
	Session     Session   `gorm:"references:UUID"`
}

type HTTPSpan struct {
	gorm.Model
	URL      string `gorm:"not null"`
	Method   string `gorm:"not null"`
	Data     string
	Headers  string `gorm:"not null"`
	Status   int    `gorm:"not null"`
	Response string

	SessionUUID uuid.UUID `gorm:"not null"`
	Session     Session   `gorm:"references:UUID"`
}

type Event struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Value string `gorm:"not null"`

	SessionUUID uuid.UUID `gorm:"not null"`
	Session     Session   `gorm:"references:UUID"`
}

type PerformanceSpan struct {
	gorm.Model
	Name  string  `gorm:"not null"`
	Value float64 `gorm:"not null"`
	URL   string  `gorm:"not null"`

	SessionUUID uuid.UUID `gorm:"not null"`
	Session     Session   `gorm:"references:UUID"`
}

type User struct {
	gorm.Model
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Role     int    `gorm:"not null"`
}

const (
	LCP = "LCP"
	CLS = "CLS"
	FID = "FID"
)
