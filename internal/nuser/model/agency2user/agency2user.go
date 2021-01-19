package agency2user

import (
	"context"
	"github.com/myxy99/ndisk/internal/nuser/model"
	"gorm.io/gorm"
	"time"
)

type AgencyUser struct {
	ID        uint `gorm:"primarykey"`
	UserId    uint `gorm:"user_id;primaryKey"`
	AgencyId  uint `gorm:"agency_id;primaryKey"`
	Status    uint `gorm:"status;default:1"` //1为正常 2为拒绝 3 为等待接受邀请
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *AgencyUser) TableName() string {
	return "agency_user"
}

func (m *AgencyUser) Add(ctx context.Context) error {
	return model.MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *AgencyUser) Adds(ctx context.Context, data []AgencyUser) (count int64, err error) {
	tx := model.MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *AgencyUser) Del(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Delete(m)
	err = tx.Error
	count = tx.RowsAffected
	return
}
func (m *AgencyUser) GetAll(ctx context.Context, data *[]AgencyUser, wheres map[string][]interface{}) (err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	err = db.Find(&data).Error
	return
}
func (m *AgencyUser) Get(ctx context.Context, start int, size int, data *[]AgencyUser, wheres map[string][]interface{}, isDelete bool) (total int64, err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	tx := db.Limit(size).Offset(start).Find(data)
	total = tx.RowsAffected
	err = tx.Error
	return
}

func (m *AgencyUser) GetById(ctx context.Context, IgnoreDel bool) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	return db.First(m).Error
}

func (m *AgencyUser) GetByWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.First(m).Error
}

func (m *AgencyUser) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *AgencyUser) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *AgencyUser) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m *AgencyUser) UpdateStatus(ctx context.Context, status uint32) error {
	return model.MainDB().Table(m.TableName()).WithContext(ctx).Where("id=?", m.ID).Update("status", status).Error
}

func (m *AgencyUser) DelRes(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Update("deleted_at", nil)
	err = tx.Error
	count = tx.RowsAffected
	return
}
