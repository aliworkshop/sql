package sql

import (
	"context"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"gorm.io/gorm"
)

func (db *db) getTx(query dbcore.QueryModel) *gorm.DB {
	if query == nil {
		return nil
	}
	iTx := query.GetTransaction()
	if iTx != nil {
		tx := iTx.(*gorm.DB)
		return tx
	}
	return nil
}

func (db *db) StartTransaction(query dbcore.QueryModel) (err error.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		tx = db.gormDB.Begin()
	}
	if tx.Error != nil {
		err = error.Internal(tx.Error)
		return
	}
	query.SetTransaction(tx)
	return
}

func (db *db) CommitTransaction(query dbcore.QueryModel) (err error.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		return
	}
	dbc := tx.Commit()
	if dbc.Error != nil {
		err = error.Internal(dbc.Error)
		return
	}
	return
}

func (db *db) RollbackTransaction(query dbcore.QueryModel) (err error.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		return
	}
	dbc := tx.Rollback()
	if dbc.Error != nil {
		err = error.Internal(dbc.Error)
		return
	}
	return
}

func (db *db) FinalizeTransaction(ctx context.Context, query dbcore.QueryModel,
	err error.ErrorModel) error.ErrorModel {
	if err == nil {
		err = error.HandleError(ctx.Err())
	}
	if err != nil {
		e := db.RollbackTransaction(query)
		if rErr := error.HandleError(e); rErr != nil {
			return rErr
		}
		return err
	}
	err = db.CommitTransaction(query)
	if err = error.HandleError(err); err != nil {
		return err
	}
	return nil
}
