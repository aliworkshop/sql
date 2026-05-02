package sql

import (
	"github.com/aliworkshop/errors"
)

func errorHandler(err errors.ErrorModel) errors.ErrorModel {
	err = err.WithSource("sql")
	return err
}
