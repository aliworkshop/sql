package sql

import (
	"context"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/error"
	"gorm.io/gorm"
)

func (db *repo) getTx(query dbcore.QueryModel) *gorm.DB {
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

func (db *repo) StartTransaction(query dbcore.QueryModel) (err error.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		tx = db.gormDB.Begin()
		if query.GetBody() == nil {
			tx = tx.Model(query.GetModel())
		}
	}
	if tx.Error != nil {
		err = error.Internal(tx.Error)
		return
	}
	query.SetTransaction(tx)
	return
}

func (db *repo) CommitTransaction(query dbcore.QueryModel) (err error.ErrorModel) {
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

func (db *repo) RollbackTransaction(query dbcore.QueryModel) (err error.ErrorModel) {
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

func (db *repo) FinalizeTransaction(ctx context.Context, query dbcore.QueryModel,
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
