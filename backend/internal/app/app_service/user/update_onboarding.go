package user

import (
	"context"
	"saas-template/config"
	"saas-template/generated/oapi"
	"saas-template/internal/app/mutation"
	"saas-template/internal/webserver/middleware"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

func UpdateOnboarding(
	ctx context.Context,
	app *config.App,
	onboardingCompleted bool,
) (oapi.PatchApiV1UserOnboardingResponseObject, error) {
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

	user, err := mutation.UpdateUserOnboarding(ctx, app.DB(), userID, tenantID, onboardingCompleted)
	if err != nil {
		return nil, err
	}

	return oapi.PatchApiV1UserOnboarding200JSONResponse{
		Id:                  user.ID,
		Email:               types.Email(user.Email),
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Role:                oapi.BaseUserRole(user.Role),
		OnboardingCompleted: user.OnboardingCompleted,
	}, nil
}
