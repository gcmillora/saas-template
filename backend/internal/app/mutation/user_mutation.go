package mutation

import (
	"adobo/generated/db/database/public/model"
	"adobo/generated/db/database/public/table"
	"context"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func CreateUser(
	ctx context.Context,
	db qrm.DB,
	user model.UserTbl,
) (*model.UserTbl, error) {
	ctbl := table.UserTbl

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	stmt := ctbl.INSERT(ctbl.MutableColumns).
		MODEL(user).
		RETURNING(ctbl.AllColumns)

	dest := []model.UserTbl{}
	err := stmt.QueryContext(ctx, db, &dest)

	if err != nil {
		return nil, err
	}

	insertedUser := dest[0]

	return &insertedUser, nil
}

func DeleteUser(
	ctx context.Context,
	db qrm.DB,
	userID uuid.UUID,
	tenantID uuid.UUID,
) error {
	utbl := table.UserTbl

	stmt := utbl.DELETE().
		WHERE(
			pg.AND(
				utbl.ID.EQ(pg.UUID(userID)),
				utbl.TenantID.EQ(pg.UUID(tenantID)),
			),
		)

	_, err := stmt.ExecContext(ctx, db)

	return err
}