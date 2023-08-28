package sql

import (
	"context"
	"fmt"
	"github.com/aliworkshop/configer"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type db struct {
	config      config
	queryParser dbcore.QueryParser
	gormDB      *gorm.DB
}

func NewRepository(configRegistry configer.Registry, parser dbcore.QueryParser) dbcore.RDBMS {
	db := new(db)
	// load config
	err := configRegistry.Root().Unmarshal(&db.config)
	if err != nil {
		panic(err)
	}
	err = configRegistry.Unmarshal(&db.config.Sql)
	if err != nil {
		panic(err)
	}
	//
	db.queryParser = parser
	return db
}

func (db *db) DB() interface{} {
	return db.gormDB
}

func (db *db) GetGormDB(queries ...dbcore.QueryModel) *gorm.DB {
	var q dbcore.QueryModel
	if queries != nil && len(queries) > 0 {
		q = queries[0]
	}
	gormDb := db.gormDB
	if q != nil {
		if db := q.GetDB(); db != nil {
			return db.(*gorm.DB)
		}
		tr := db.GetTransaction(q)
		if tr != nil {
			return tr.(*gorm.DB)
		}
		body := q.GetBody()
		if body != nil {
			switch body.(type) {
			case map[string]interface{}:
				break
			default:
				return gormDb.Model(body)
			}
		}
		model := q.GetModel()
		if model != nil {
			return gormDb.Model(model)
		}
	}
	return gormDb
}

func (db *db) GetDB(queries ...dbcore.QueryModel) interface{} {
	return db.GetGormDB(queries...)
}

func (db *db) GetTransaction(query dbcore.QueryModel) (transaction interface{}) {
	return query.GetTransaction()
}

func (db *db) Initialize() error.ErrorModel {
	if db.config.Sql.Dialect == "" {
		panic("sql dialect is not determined")
	}
	connStr := db.config.Sql.ConnectionString
	if connStr == "" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			db.config.Sql.Host,
			db.config.Sql.Port,
			db.config.Sql.Username,
			db.config.Sql.DbName,
			db.config.Sql.Password,
		)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             5 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn,     // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)
	var dialect gorm.Dialector
	switch db.config.Sql.Dialect {
	case "mysql":
		dialect = mysql.Open(connStr)
	case "postgres":
		dialect = postgres.New(postgres.Config{
			DSN:                  connStr,
			PreferSimpleProtocol: true,
		})
	}
	d, err := gorm.Open(dialect, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return error.Internal(err)
	}
	sqlDB, err := d.DB()
	if err != nil {
		return error.Internal(err)
	}
	if db.config.Sql.MaxIdleConnections != nil {
		sqlDB.SetMaxIdleConns(*db.config.Sql.MaxIdleConnections)
	}
	if db.config.Sql.MaxOpenConnections != nil {
		sqlDB.SetMaxOpenConns(*db.config.Sql.MaxOpenConnections)
	}
	if db.config.Sql.MaxLifetimeSeconds != nil {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(*db.config.Sql.MaxLifetimeSeconds))
	}
	if db.config.Debug {
		fmt.Println("debug is true")
		d = d.Debug()
	}
	db.gormDB = d
	return nil
}

func (db *db) Ping(ctx context.Context) error.ErrorModel {
	if db.gormDB == nil {
		return error.DefaultValidationError.WithDetail("db not initialized")
	}
	sqlDb, err := db.gormDB.DB()
	if err != nil {
		return error.HandleError(err)
	}
	err = sqlDb.PingContext(ctx)
	if err != nil {
		return error.Internal(err)
	}
	return nil
}
