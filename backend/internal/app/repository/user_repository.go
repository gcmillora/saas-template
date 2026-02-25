package repository

import (
	"adobo/generated/db/database/public/model"
	"adobo/generated/db/database/public/table"
	"context"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func GetUserByID(
	ctx context.Context,
	db qrm.DB,
	userID uuid.UUID,
	tenantID uuid.UUID,
) (*model.UserTbl, error) {
	ctbl := table.UserTbl

	stmt := pg.SELECT(ctbl.AllColumns).
		FROM(ctbl).
		WHERE(
			pg.AND(
				ctbl.ID.EQ(pg.UUID(userID)),
				ctbl.TenantID.EQ(pg.UUID(tenantID)),
			),
		)

	dest := model.UserTbl{}
	err := stmt.QueryContext(ctx, db, &dest)

	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func GetUserByEmail(
	ctx context.Context,
	db qrm.DB,
	email string,
) (*model.UserTbl, error) {
	ctbl := table.UserTbl

	stmt := pg.SELECT(ctbl.AllColumns).
		FROM(ctbl).
		WHERE(
			ctbl.Email.EQ(pg.String(email)),
		)

	dest := model.UserTbl{}
	err := stmt.QueryContext(ctx, db, &dest)

	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func GetUsers(
	ctx context.Context,
	db qrm.DB,
	tenantID uuid.UUID,
) (*[]model.UserTbl, error) {
	ctbl := table.UserTbl

	stmt := pg.SELECT(ctbl.AllColumns).
		FROM(ctbl).
		WHERE(
			pg.AND(
				ctbl.TenantID.EQ(pg.UUID(tenantID)),
			),
		)

	rows := []model.UserTbl{}
	err := stmt.QueryContext(ctx, db, &rows)

	if err != nil {
		return nil, err
	}

	return &rows, nil
}

