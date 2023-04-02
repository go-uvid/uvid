package daos

import (
	"luvsic3/uvid/dtos"
	"luvsic3/uvid/models"

	"github.com/google/uuid"
)

func (dao *Dao) CreateSession(dto *dtos.SessionDTO) uuid.UUID {
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
	dao.db.Create(&model)
	return model.UUID
}
