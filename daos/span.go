package daos

import (
	"github.com/go-uvid/uvid/dtos"
	"github.com/go-uvid/uvid/models"

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
	err := dao.DB.Create(&model).Error
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
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateEvent(sessionUUID uuid.UUID, dto *dtos.EventDTO) (*models.Event, error) {
	model := models.Event{
		SessionUUID: sessionUUID,
		Action:      dto.Action,
		Value:       dto.Value,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreatePerformance(sessionUUID uuid.UUID, dto *dtos.PerformanceDTO) (*models.Performance, error) {
	model := models.Performance{
		SessionUUID: sessionUUID,
		Name:        dto.Name,
		Value:       dto.Value,
		URL:         dto.URL,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateHTTP(sessionUUID uuid.UUID, dto *dtos.HTTPDTO) (*models.HTTP, error) {
	model := models.HTTP{
		SessionUUID: sessionUUID,
		Resource:    dto.Resource,
		Method:      dto.Method,
		Headers:     dto.Headers,
		Status:      dto.Status,
		Body:        dto.Body,
		Response:    dto.Response,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreatePageView(sessionUUID uuid.UUID, dto *dtos.PageViewDTO) (*models.PageView, error) {
	model := models.PageView{
		SessionUUID: sessionUUID,
		URL:         dto.URL,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}
