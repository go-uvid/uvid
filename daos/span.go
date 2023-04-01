package daos

import (
	"luvsic3/uvid/models"
)

func (dao *Dao) SaveJSError(model models.JSError) {
	dao.db.Create(model)
}

func (dao *Dao) SaveEvent(model models.Event) {
	dao.db.Create(model)
}

func (dao *Dao) SavePerformance(model models.PerformanceSpan) {
	dao.db.Create(model)
}

func (dao *Dao) SaveHTTP(model models.HTTPSpan) {
	dao.db.Create(model)
}
