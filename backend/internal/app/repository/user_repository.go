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
	customerID uuid.UUID,
	tenantID uuid.UUID,
) (*[]model.UserTbl, error) {
	ctbl := table.UserTbl

	stmt := pg.SELECT(ctbl.AllColumns).
		FROM(ctbl).
		WHERE(
			pg.AND(
				ctbl.ID.EQ(pg.UUID(customerID)),
			),
		)
	
	rows := []model.UserTbl{}
	err := stmt.QueryContext(ctx, db, &rows)

	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func GetUserByAuthID(
	ctx context.Context,
	db qrm.DB,
	authID uuid.UUID,
) (*model.UserTbl, error ){
	ctbl := table.UserTbl

	stmt := pg.SELECT(ctbl.AllColumns).
		FROM(ctbl).
		WHERE(
			pg.AND(
				ctbl.AuthID.EQ(pg.UUID(authID)),
			),
		)
	
	rows := []model.UserTbl{}
	err := stmt.QueryContext(ctx, db, &rows)

	if err != nil {
		return nil, err
	}

	return &rows[0], nil
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

func GetUsersPaginated(
	ctx context.Context,
	db qrm.DB,
	tenantID uuid.UUID,
	params PaginationParams,
) (*PaginatedResponse[model.CategoryTbl], error) {
	ctbl := table.CategoryTbl

	whereConditions := pg.AND(
		ctbl.TenantID.EQ(pg.UUID(tenantID)),
	)

	countStmt := pg.SELECT(pg.COUNT(ctbl.ID)).
		FROM(ctbl).
		WHERE(whereConditions)
	
	var totalCount struct {
		Count int64
	}

	err := countStmt.QueryContext(ctx, db, &totalCount)
	if err != nil {
		return nil, err
	}

	baseQuery := pg.SELECT(ctbl.AllColumns).FROM(ctbl).WHERE(whereConditions)
	stmt := params.ApplyPagination(baseQuery)

	rows := []model.CategoryTbl{}
	err = stmt.QueryContext(ctx, db, &rows)

	if err != nil {
		return nil, err
	}

	return NewPaginatedResponse(rows, int(totalCount.Count), params), nil
}