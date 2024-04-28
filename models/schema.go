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

	JSErrors     []JSError
	HTTPs        []HTTP
	Events       []Event
	Performances []Performance
	PageViews    []PageView
}

type JSError struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Message string `gorm:"not null"`
	Stack   string `gorm:"not null"`
	Cause   string

	SessionID uint `gorm:"not null"`
}

type HTTP struct {
	gorm.Model
	Resource string `gorm:"not null"`
	Method   string `gorm:"not null"`
	Headers  string `gorm:"not null"`
	Status   int    `gorm:"not null"`
	Body     string
	Response string

	SessionID uint `gorm:"not null"`
}

type Event struct {
	gorm.Model
	Action string `gorm:"not null"`
	Value  string `gorm:"not null"`

	SessionID uint `gorm:"not null"`
}

type Performance struct {
	gorm.Model
	Name  string  `gorm:"not null"`
	Value float64 `gorm:"not null"`
	URL   string  `gorm:"not null"`

	SessionID uint `gorm:"not null"`
}

type PageView struct {
	gorm.Model
	URL string `gorm:"not null"`

	SessionID uint `gorm:"not null"`
}

const (
	LCP = "LCP"
	CLS = "CLS"
	FID = "FID"
)
