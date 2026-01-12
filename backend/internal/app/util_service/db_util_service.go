package utils

import (
	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func UUIDs(uuids []uuid.UUID) []pg.Expression {
	result := make([]pg.Expression, len(uuids))

	for i, uuid := range uuids {
		result[i] = pg.UUID(uuid)
	}

	return result
}

