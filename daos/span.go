package daos

import (
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/models"

	"github.com/google/uuid"
)

func (dao *Dao) CreateSession(dto *dtos.SessionDTO) (*models.Session, error) {
	model := models.Session{
		UA:         dto.UA,
		Language:   dto.Language,
		IP:         dto.IP,
		AppVersion: dto.AppVersion,
		URL:        dto.URL,
		Screen:     dto.Screen,
		Referrer:   dto.Referrer,
		Meta:       dto.Meta,
		UUID:       uuid.New(),
	}
	err := dao.db.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateJSError(sessionUUID uuid.UUID, dto *dtos.ErrorDTO) (*models.JSError, error) {
	model := models.JSError{
		SessionUUID: sessionUUID,
		Name:        dto.Name,
		Message:     dto.Message,
		Stack:       dto.Stack,
		Cause:       dto.Cause,
	}
	err := dao.db.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateEvent(sessionUUID uuid.UUID, dto *dtos.EventDTO) (*models.Event, error) {
	model := models.Event{
		SessionUUID: sessionUUID,
		Name:        dto.Name,
		Value:       dto.Value,
	}
	err := dao.db.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreatePerformance(sessionUUID uuid.UUID, dto *dtos.PerformanceDTO) (*models.PerformanceSpan, error) {
	model := models.PerformanceSpan{
		SessionUUID: sessionUUID,
		Name:        dto.Name,
		Value:       dto.Value,
		URL:         dto.URL,
	}
	err := dao.db.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateHTTP(sessionUUID uuid.UUID, dto *dtos.HTTPDTO) (*models.HTTPSpan, error) {
	model := models.HTTPSpan{
		SessionUUID: sessionUUID,
		URL:         dto.URL,
		Method:      dto.Method,
		Headers:     dto.Headers,
		Status:      dto.Status,
		Data:        dto.Data,
		Response:    dto.Response,
	}
	err := dao.db.Create(&model).Error
	return &model, err
}
