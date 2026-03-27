package handler

import (
	"saas-template/config"
	"saas-template/generated/oapi"
	oapi_admin "saas-template/generated/oapi/admin"
	oapi_public "saas-template/generated/oapi/public"
)

type CombinedStrictServerInterface interface {
	oapi.StrictServerInterface
	oapi_public.StrictServerInterface
	oapi_admin.StrictServerInterface
}

type Handler struct {
	app *config.App
	CombinedStrictServerInterface
}

func NewHandler(app *config.App) *Handler {
	handler := Handler{
		app: app,
	}

	return &handler
}
