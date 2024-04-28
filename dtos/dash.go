package dtos

import (
	"time"

	"github.com/go-uvid/uvid/tools"
)

type (
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
