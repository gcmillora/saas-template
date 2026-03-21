package handler

import (
	"context"
	"saas-template/generated/oapi"
	"saas-template/internal/app/app_service/user"
)

// GET (/user)
func (h *Handler) GetApiV1User(
	ctx context.Context,
	request oapi.GetApiV1UserRequestObject,
) (oapi.GetApiV1UserResponseObject, error) {
	data, err := user.GetUser(ctx, h.app)
	if err != nil {
		return nil, err
	}

	return data, nil
}
