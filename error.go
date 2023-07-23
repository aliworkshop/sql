package sql

import (
	"github.com/aliworkshop/error"
)

func errorHandler(err error.ErrorModel) error.ErrorModel {
	err = err.WithSource("sql")
	return err
}
