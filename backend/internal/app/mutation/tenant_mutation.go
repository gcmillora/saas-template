package mutation

import (
	"saas-template/generated/db/database/public/model"
	"saas-template/generated/db/database/public/table"
	"context"
	"time"

	"github.com/go-jet/jet/v2/qrm"
)

func CreateTenant(
	ctx context.Context,
	db qrm.DB,
	tenant model.TenantTbl,
) (*model.TenantTbl, error) {
	ctbl := table.TenantTbl

	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	stmt := ctbl.INSERT(ctbl.MutableColumns).
		MODEL(tenant).
		RETURNING(ctbl.AllColumns)

	dest := []model.TenantTbl{}
	err := stmt.QueryContext(ctx, db, &dest)

	if err != nil {
		return nil, err
	}

	return &dest[0], nil
}
