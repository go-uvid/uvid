package dtos

import "time"

type (
	LoginDTO struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	UpdatePasswordDTO struct {
		Password string `json:"password" validate:"required"`
	}
	TimeRangeDTO struct {
		Start time.Time `json:"start" validate:"required"`
		End   time.Time `json:"end" validate:"required"`
		Unit  string    `json:"unit" validate:"required"`
	}
)

const (
	UnitHour  = "hour"
	UnitDay   = "day"
	UnitWeek  = "week"
	UnitMonth = "month"
)
