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
		Start time.Time `query:"start" validate:"required"`
		End   time.Time `query:"end" validate:"required"`
		Unit  string    `query:"unit" validate:"required"`
	}

	IntervalData struct {
		X string `json:"x" validate:"required"`
		Y int64  `json:"y" validate:"required"`
	}
)

const (
	UnitHour  = "hour"
	UnitDay   = "day"
	UnitWeek  = "week"
	UnitMonth = "month"
)
