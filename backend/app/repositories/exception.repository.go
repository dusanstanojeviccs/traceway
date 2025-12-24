package repositories

import (
	"backend/app/models"
	"database/sql"

	"github.com/tracewayapp/go-lightning/lpg"
)

type exceptionRepository struct{}

func (e exceptionRepository) FindAll(tx *sql.Tx) ([]*models.Exception, error) {
	return lpg.SelectGeneric[models.Exception](tx, "SELECT * FROM exceptions WHERE archived = false")
}

var ExceptionRepository = exceptionRepository{}
