package authentication

import (
	"context"
	"saas-template/config"
	"saas-template/internal/app/app_service/audit"

	"github.com/google/uuid"
)

func SignOut(ctx context.Context, app *config.App, userID uuid.UUID, tenantID uuid.UUID) {
	audit.Log(ctx, app, "signout", &userID, &tenantID, nil)
}
