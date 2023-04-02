package daos

import (
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/models"

	"github.com/google/uuid"
)

func (dao *Dao) CreateJSError(sessionUUID uuid.UUID, dto *dtos.ErrorDTO) {
	model := models.JSError{
		SessionUUID: sessionUUID,
		Name:        dto.Name,
		Message:     dto.Message,
		Stack:       dto.Stack,
		Cause:       dto.Cause,
	}
	dao.db.Create(&model)
}

func (dao *Dao) CreateEvent(model *models.Event) {
	dao.db.Create(&model)
}

func (dao *Dao) CreatePerformance(model *models.PerformanceSpan) {
	dao.db.Create(&model)
}

func (dao *Dao) CreateHTTP(model *models.HTTPSpan) {
	dao.db.Create(&model)
}
