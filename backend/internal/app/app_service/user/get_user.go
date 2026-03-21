package user

import (
	"saas-template/config"
	"saas-template/generated/oapi"
	"saas-template/internal/app/repository"
	"saas-template/internal/webserver/middleware"
	"context"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

type BaseUser = oapi.BaseUser

func GetUser(
	ctx context.Context,
	app *config.App,
) (oapi.GetApiV1UserResponseObject, error) {
	sessionData, err := middleware.GetSessionData(ctx, app.Session())
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(sessionData.UserID)
	if err != nil {
		return nil, err
	}

	tenantID, err := uuid.Parse(sessionData.TenantID)
	if err != nil {
		return nil, err
	}

	user, err := repository.GetUserByID(ctx, app.DB(), userID, tenantID)
	if err != nil {
		return nil, err
	}

	return oapi.GetApiV1User200JSONResponse{
		Email:     types.Email(user.Email),
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
