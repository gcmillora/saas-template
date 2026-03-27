package authentication

import (
	"context"
	"errors"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/app_service/audit"
	"saas-template/internal/app/mutation"
	"saas-template/internal/app/repository"
	"saas-template/internal/app/util_service/password"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type PostSignupBody struct {
	Email           string
	Password        string
	ConfirmPassword string
	FirstName       *string
	LastName        *string
}

func SignUp(ctx context.Context, app *config.App, body PostSignupBody) (*model.UserTbl, error) {
	if body.Password != body.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	if errs := password.ValidateComplexity(body.Password); len(errs) > 0 {
		return nil, errors.New("Password does not meet requirements: " + strings.Join(errs, ", "))
	}

	existing, _ := repository.GetUserByEmail(ctx, app.DB(), body.Email)
	if existing != nil {
		return nil, errors.New("unable to create account")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashStr := string(hash)

	tenant, err := mutation.CreateTenant(ctx, app.DB(), model.TenantTbl{
		Name: body.Email,
	})
	if err != nil {
		return nil, err
	}

	user := model.UserTbl{
		Email:        body.Email,
		PasswordHash: &hashStr,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		AuthProvider: "email",
		TenantID:     tenant.ID,
		Role:         "user",
	}

	created, err := mutation.CreateUser(ctx, app.DB(), user)
	if err != nil {
		return nil, err
	}

	audit.Log(ctx, app, "signup", &created.ID, &created.TenantID, nil)
	return created, nil
}
