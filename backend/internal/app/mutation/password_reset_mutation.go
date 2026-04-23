package mutation

import (
	"context"
	"saas-template/generated/db/database/public/model"
	"saas-template/generated/db/database/public/table"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func CreatePasswordReset(
	ctx context.Context,
	db qrm.DB,
	reset model.PasswordResetTbl,
) error {
	ctbl := table.PasswordResetTbl

	reset.CreatedAt = time.Now()

	stmt := ctbl.INSERT(ctbl.MutableColumns).
		MODEL(reset)

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func InvalidateAllUserTokens(
	ctx context.Context,
	db qrm.DB,
	userID uuid.UUID,
) error {
	ctbl := table.PasswordResetTbl
	now := time.Now()

	stmt := ctbl.UPDATE(ctbl.UsedAt).
		SET(now).
		WHERE(
			pg.AND(
				ctbl.UserID.EQ(pg.UUID(userID)),
				ctbl.UsedAt.IS_NULL(),
			),
		)

	_, err := stmt.ExecContext(ctx, db)
	return err
}
