package authentication

import (
	"context"
	"errors"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/app_service/audit"
	"saas-template/internal/app/repository"

	"golang.org/x/crypto/bcrypt"
)

func SignIn(
	ctx context.Context,
	app *config.App,
	email string,
	passwordInput string,
) (*model.UserTbl, error) {
	user, err := repository.GetUserByEmail(ctx, app.DB(), email)
	if err != nil {
		audit.Log(ctx, app, "signin_failed", nil, nil, map[string]string{"email": email})
		return nil, errors.New("invalid email or password")
	}

	if user.PasswordHash == nil {
		audit.Log(
			ctx,
			app,
			"signin_failed",
			&user.ID,
			&user.TenantID,
			map[string]string{"email": email},
		)
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.PasswordHash),
		[]byte(passwordInput),
	); err != nil {
		audit.Log(
			ctx,
			app,
			"signin_failed",
			&user.ID,
			&user.TenantID,
			map[string]string{"email": email},
		)
		return nil, errors.New("invalid email or password")
	}

	audit.Log(ctx, app, "signin_success", &user.ID, &user.TenantID, nil)
	return user, nil
}
