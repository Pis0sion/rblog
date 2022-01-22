package mysql

import (
	"fmt"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	"github.com/Pis0sion/rblog/internal/rblog/opts"
	"github.com/Pis0sion/rblog/pkg/db"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	dnIns *gorm.DB
}

func (d datastore) Articles() dto.ArticleDto {
	return newArticles(d.dnIns)
}

func (d datastore) ArticleTags() dto.ArticleTagsDto {
	return newArticleTags(d.dnIns)
}

var (
	myFactory dto.Factory
	once      sync.Once
)

// GetDatabaseFactoryEntity initialize mysql
// get mysql factory
func GetDatabaseFactoryEntity(opts *opts.MysqlOpts) (dto.Factory, error) {

	if opts == nil && myFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB

	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			Charset:               opts.Charset,
			ParseTime:             opts.ParseTime,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              0,
			Logger:                nil,
		}
		dbIns, err = db.New(options)
		myFactory = &datastore{dbIns}
	})

	if myFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", myFactory, err)
	}

	return myFactory, nil
}
