package handler

import (
	"context"
	oapi_admin "saas-template/generated/oapi/admin"
)

func (h *Handler) GetApiAdminV1Health(
	ctx context.Context,
	request oapi_admin.GetApiAdminV1HealthRequestObject,
) (oapi_admin.GetApiAdminV1HealthResponseObject, error) {
	return oapi_admin.GetApiAdminV1Health200JSONResponse{
		Message: ptr("ok"),
	}, nil
}

func ptr(s string) *string { return &s }
