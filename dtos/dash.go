package dtos

import (
	"luvsic3/uvid/tools"
	"time"
)

type (
	LoginDTO struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	UpdatePasswordDTO struct {
		Password string `json:"password" validate:"required"`
	}
	SpanFilterDTO struct {
		Start time.Time `query:"start" validate:"required"`
		End   time.Time `query:"end" validate:"required"`
	}
	TimeIntervalSpanDTO struct {
		SpanFilterDTO
		Unit tools.Unit `query:"unit" validate:"required"`
	}

	IntervalData struct {
		X string  `json:"x" validate:"required"`
		Y float64 `json:"y" validate:"required"`
	}
)
