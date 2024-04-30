package daos

import (
	"github.com/rick-you/uvid/dtos"
	"github.com/rick-you/uvid/models"
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
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateJSError(sessionID uint, dto *dtos.ErrorDTO) (*models.JSError, error) {
	model := models.JSError{
		SessionID: sessionID,
		Name:      dto.Name,
		Message:   dto.Message,
		Stack:     dto.Stack,
		Cause:     dto.Cause,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateEvent(sessionID uint, dto *dtos.EventDTO) (*models.Event, error) {
	model := models.Event{
		SessionID: sessionID,
		Action:    dto.Action,
		Value:     dto.Value,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreatePerformance(sessionID uint, dto *dtos.PerformanceDTO) (*models.Performance, error) {
	model := models.Performance{
		SessionID: sessionID,
		Name:      dto.Name,
		Value:     dto.Value,
		URL:       dto.URL,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreateHTTP(sessionID uint, dto *dtos.HTTPDTO) (*models.HTTP, error) {
	model := models.HTTP{
		SessionID: sessionID,
		Resource:  dto.Resource,
		Method:    dto.Method,
		Headers:   dto.Headers,
		Status:    dto.Status,
		Body:      dto.Body,
		Response:  dto.Response,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}

func (dao *Dao) CreatePageView(sessionID uint, dto *dtos.PageViewDTO) (*models.PageView, error) {
	model := models.PageView{
		SessionID: sessionID,
		URL:       dto.URL,
	}
	err := dao.DB.Create(&model).Error
	return &model, err
}
