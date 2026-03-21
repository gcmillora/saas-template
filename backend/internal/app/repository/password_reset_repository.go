package repository

import (
	"saas-template/generated/db/database/public/model"
	"saas-template/generated/db/database/public/table"
	"context"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func GetPasswordResetByTokenHash(
	ctx context.Context,
	db qrm.DB,
	tokenHash string,
) (*model.PasswordResetTbl, error) {
	ctbl := table.PasswordResetTbl

	stmt := pg.SELECT(ctbl.AllColumns).
		FROM(ctbl).
		WHERE(
			pg.AND(
				ctbl.TokenHash.EQ(pg.String(tokenHash)),
				ctbl.UsedAt.IS_NULL(),
				ctbl.ExpiresAt.GT(pg.TimestampT(time.Now())),
			),
		)

	dest := model.PasswordResetTbl{}
	err := stmt.QueryContext(ctx, db, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}
