package sql

import (
	"context"
	"github.com/aliworkshop/dbcore"
	"github.com/aliworkshop/errors"
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

func (db *repo) StartTransaction(query dbcore.QueryModel) (err errors.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		tx = db.gormDB.Begin()
	}
	if tx.Error != nil {
		err = errors.Internal(tx.Error)
		return
	}
	query.SetTransaction(tx)
	return
}

func (db *repo) CommitTransaction(query dbcore.QueryModel) (err errors.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		return
	}
	dbc := tx.Commit()
	if dbc.Error != nil {
		err = errors.Internal(dbc.Error)
		return
	}
	return
}

func (db *repo) RollbackTransaction(query dbcore.QueryModel) (err errors.ErrorModel) {
	tx := db.getTx(query)
	if tx == nil {
		return
	}
	dbc := tx.Rollback()
	if dbc.Error != nil {
		err = errors.Internal(dbc.Error)
		return
	}
	return
}

func (db *repo) FinalizeTransaction(ctx context.Context, query dbcore.QueryModel,
	err errors.ErrorModel) errors.ErrorModel {
	if err != nil {
		e := db.RollbackTransaction(query)
		if rErr := errors.HandleError(e); rErr != nil {
			return rErr
		}
		return err
	}
	err = db.CommitTransaction(query)
	if err = errors.HandleError(err); err != nil {
		return err
	}
	return nil
}

func (db *repo) RunInTransaction(ctx context.Context, fn func(dbcore.QueryModel) errors.ErrorModel, label string) errors.ErrorModel {
	q := dbcore.NewQuery().WithContext(ctx)
	if err := db.StartTransaction(q); err != nil {
		db.logger.Error("start tx", "op", label, "err", err)
		return err
	}
	var txErr errors.ErrorModel
	defer func() {
		if e := db.FinalizeTransaction(context.Background(), q, txErr); e != nil {
			db.logger.Error("finalize tx", "op", label, "err", e)
		}
	}()
	if txErr = fn(q); txErr != nil {
		db.logger.Error("tx step failed", "op", label, "err", txErr)
		return txErr
	}
	return nil
}
