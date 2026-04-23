package repository

import (
	"context"
	"saas-template/generated/db/database/public/model"
	"saas-template/generated/db/database/public/table"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func GetTenantByID(
	ctx context.Context,
	db qrm.DB,
	tenantID uuid.UUID,
) (*[]model.TenantTbl, error) {
	ttbl := table.TenantTbl

	stmt := pg.SELECT(ttbl.AllColumns).
		FROM(ttbl).
		WHERE(
			pg.AND(
				ttbl.ID.EQ(pg.UUID(tenantID)),
			),
		)

	rows := []model.TenantTbl{}
	err := stmt.QueryContext(ctx, db, &rows)

	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func GetTenants(
	ctx context.Context,
	db qrm.DB,
) (*[]model.TenantTbl, error) {
	ttbl := table.TenantTbl

	stmt := pg.SELECT(ttbl.AllColumns).
		FROM(ttbl)

	rows := []model.TenantTbl{}
	err := stmt.QueryContext(ctx, db, &rows)

	if err != nil {
		return nil, err
	}

	return &rows, nil
}
