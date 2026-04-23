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

func UpdateUserPassword(
	ctx context.Context,
	db qrm.DB,
	userID uuid.UUID,
	passwordHash string,
) error {
	utbl := table.UserTbl
	now := time.Now()

	stmt := utbl.UPDATE(utbl.PasswordHash, utbl.UpdatedAt).
		SET(passwordHash, now).
		WHERE(utbl.ID.EQ(pg.UUID(userID)))

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func UpdateUserOnboarding(
	ctx context.Context,
	db qrm.DB,
	userID uuid.UUID,
	tenantID uuid.UUID,
	onboardingCompleted bool,
) (model.UserTbl, error) {
	utbl := table.UserTbl
	now := time.Now()

	var dest model.UserTbl
	stmt := utbl.UPDATE(
		utbl.OnboardingCompleted,
		utbl.UpdatedAt,
	).SET(
		onboardingCompleted,
		now,
	).WHERE(
		pg.AND(
			utbl.ID.EQ(pg.UUID(userID)),
			utbl.TenantID.EQ(pg.UUID(tenantID)),
		),
	).RETURNING(utbl.AllColumns)

	err := stmt.QueryContext(ctx, db, &dest)
	return dest, err
}
