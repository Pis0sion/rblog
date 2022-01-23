package mod

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Title     string    `gorm:"column:title"`
	State     int       `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:createdAt"`
}

func (a Article) Create(db *gorm.DB) error {

	return db.Create(&a).Error
}

func (a Article) List(db *gorm.DB, pageOffset, pageSize int) (articles []*Article, count int64, err error) {

	if pageOffset >= 0 && pageSize >= 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	if a.Title != "" {
		db = db.Where("title = ?", a.Title)
	}
	db = db.Where("status = ?", a.State)

	if err = db.Find(&articles).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return articles, count, err
}
