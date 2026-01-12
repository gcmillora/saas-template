package user

import (
	"adobo/config"
	"adobo/generated/oapi"
	"adobo/internal/app/repository"
	"adobo/internal/webserver/middleware"
	"context"
	"log/slog"

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

	authID, err := uuid.Parse(sessionData.UserID)
	slog.Default().Info(authID.String())
	if err != nil {
		return nil, err
	}

	data, err := repository.GetUserByAuthID(ctx, app.DB(), authID)
	if err != nil {
		return nil, err
	}

	return oapi.GetApiV1User200JSONResponse{
		Email: types.Email(*data.Email),
		Id: data.ID,
		FirstName: data.FirstName,
		LastName: data.LastName,
	}, nil
}