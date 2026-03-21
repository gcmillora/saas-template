package authentication

import (
	"context"
	"errors"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
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
		return nil, errors.New("invalid email or password")
	}

	if user.PasswordHash == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.PasswordHash),
		[]byte(passwordInput),
	); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
