package dtos

import (
	"time"

	"github.com/rick-you/uvid/tools"
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

	CountDTO struct {
		Pv        int64 `json:"pv" validate:"required"`
		Uv        int64 `json:"uv" validate:"required"`
		JsError   int64 `json:"jsError" validate:"required"`
		HttpError int64 `json:"httpError" validate:"required"`
	}
)
