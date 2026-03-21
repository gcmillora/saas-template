package handler

import (
	"context"
	"saas-template/generated/oapi"
	"saas-template/internal/app/app_service/user"
)

// PATCH (/user/onboarding)
func (h *Handler) PatchApiV1UserOnboarding(
	ctx context.Context,
	request oapi.PatchApiV1UserOnboardingRequestObject,
) (oapi.PatchApiV1UserOnboardingResponseObject, error) {
	return user.UpdateOnboarding(ctx, h.app, request.Body.OnboardingCompleted)
}
