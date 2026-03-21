package authentication

import (
	"context"
	"saas-template/config"

	"github.com/google/uuid"
)

func SignOut(ctx context.Context, app *config.App, userID uuid.UUID, tenantID uuid.UUID) {
	_ = ctx
	_ = app
}
