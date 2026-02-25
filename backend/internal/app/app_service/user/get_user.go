package user

import (
	"adobo/config"
	"adobo/generated/oapi"
	"adobo/internal/app/repository"
	"adobo/internal/webserver/middleware"
	"context"
	"errors"

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

	users, err := repository.GetUserByID(ctx, app.DB(), userID, tenantID)
	if err != nil {
		return nil, err
	}

	if len(*users) == 0 {
		return nil, errors.New("user not found")
	}

	data := (*users)[0]

	return oapi.GetApiV1User200JSONResponse{
		Email:     types.Email(*data.Email),
		Id:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}, nil
}
