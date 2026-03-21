package authentication

import (
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/mutation"
	"saas-template/internal/app/repository"
	"context"
	"errors"

	"github.com/google/uuid"
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

	if len(body.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	existing, _ := repository.GetUserByEmail(ctx, app.DB(), body.Email)
	if existing != nil {
		return nil, errors.New("a user with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashStr := string(hash)
	defaultTenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	user := model.UserTbl{
		Email:        body.Email,
		PasswordHash: &hashStr,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		AuthProvider: "email",
		TenantID:     defaultTenantID,
	}

	created, err := mutation.CreateUser(ctx, app.DB(), user)
	if err != nil {
		return nil, err
	}

	return created, nil
}
