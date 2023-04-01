package models

import (
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UA         string `gorm:"not null"`
	Language   string `gorm:"not null"`
	IP         string `gorm:"not null"`
	AppVersion string
	URL        string `gorm:"not null"`
	Screen     string `gorm:"not null"`
	Referrer   string
	Meta       string

	JSErrors    []JSError
	HTTPMetrics []HTTPMetric
	Events      []Event
	PerfMetrics []PerformanceMetric
}

type JSError struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Message string `gorm:"not null"`
	Stack   string `gorm:"not null"`
	Cause   string

	SessionID uint `gorm:"not null"`
	Session   Session
}

type HTTPMetric struct {
	gorm.Model
	URL      string `gorm:"not null"`
	Method   string `gorm:"not null"`
	Data     string
	Headers  string `gorm:"not null"`
	Status   int    `gorm:"not null"`
	Response string

	SessionID uint `gorm:"not null"`
	Session   Session
}

type Event struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Value string `gorm:"not null"`

	SessionID uint `gorm:"not null"`
	Session   Session
}

type PerformanceMetric struct {
	gorm.Model
	Name  string  `gorm:"not null"`
	Value float64 `gorm:"not null"`
	URL   string  `gorm:"not null"`

	SessionID uint `gorm:"not null"`
	Session   Session
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
