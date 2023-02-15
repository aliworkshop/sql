package sql

import (
	"github.com/aliworkshop/errorslib"
)

func errorHandler(err errorslib.ErrorModel) errorslib.ErrorModel {
	err = err.WithSource("sql")
	return err
}
