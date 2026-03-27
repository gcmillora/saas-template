package mutation

import (
	"context"
	"saas-template/generated/db/database/public/model"
	"saas-template/generated/db/database/public/table"
	"time"

	"github.com/go-jet/jet/v2/qrm"
)

func CreateAuditLog(
	ctx context.Context,
	db qrm.DB,
	auditLog model.AuditLogTbl,
) error {
	ctbl := table.AuditLogTbl

	auditLog.CreatedAt = time.Now()

	stmt := ctbl.INSERT(ctbl.MutableColumns).
		MODEL(auditLog)

	_, err := stmt.ExecContext(ctx, db)
	return err
}
